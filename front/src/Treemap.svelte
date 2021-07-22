<script>
  import { onMount } from 'svelte';
  import * as d3 from 'd3';
	import { createEventDispatcher } from 'svelte';

  export let data;
  export let height;
  export let width;

  let el;

  const hierarchyData = {
    name: ".",
    children: {}
  };

	const dispatch = createEventDispatcher();

	onMount(async () => {
    const tooltip = d3
      .select('body')
      .append('div')
      .style('position', 'absolute')
      .style('z-index', '10')
      .style('visibility', 'hidden')
      .style('background-color', 'white')
      .style('border', 'solid')
      .style('border-width', '2px')
      .style('border-radius', '5px')
      .style('padding', '5px')
      .text(`kk`);

    const format = d3.format(",d");

    const name = d => d.ancestors().reverse().map(d => d.data.name).join("/")
    const color = d3.scaleLinear()
      .domain([0, 1])
      .range(["#e1eec3", "#f05053"]);

    const dataScale = d3.scaleLinear()
        .domain([d3.min(data.files, function(d){return d.lines}), d3.max(data.files, function(d){return d.lines})])
        .range([1,10]);

    data.files.forEach((f) => {
      let path = f.directory  ? f.directory.split("/") : [];
      if (f.isDirectory) {
        path = path.concat([f.name]);
      }


      let d = hierarchyData;
      if (f.name === ".") {
        Object.keys(f).forEach(key => {
          d[key] = f[key];
        });
      }
      else {
        if (f.directory === ".") {
          path = [f.name];
        }
        path.forEach((p) => {
          if(!d.children) {
            d.children = {};
          }
          if (!d.children[p]) {
            d.children[p] = {};
          }
          d = d.children[p];
          if(!d.children) {
            d.children = {};
          }
        });

        Object.keys(f).forEach(key => {
          d[key] = f[key];
        });
      }
    });

    console.log("hierarchyData", hierarchyData)
    flattenData(hierarchyData);
    hierarchyData.lines = d3.sum(hierarchyData.children, d => d.lines);

    function tile(node, x0, y0, x1, y1) {
      d3.treemapResquarify(node, 0, 0, width, height);
      for (const child of node.children) {
        child.x0 = x0 + child.x0 / width * (x1 - x0);
        child.x1 = x0 + child.x1 / width * (x1 - x0);
        child.y0 = y0 + child.y0 / height * (y1 - y0);
        child.y1 = y0 + child.y1 / height * (y1 - y0);
      }
    }

    const treemap = data => d3.treemap()
        .tile(tile)
      (d3.hierarchy(hierarchyData)
        .sum(d => dataScale(d.lines))
        .sort((a, b) => b.value - a.value));

    const x = d3.scaleLinear().rangeRound([0, width]);
    const y = d3.scaleLinear().rangeRound([0, height]);

    const svg = d3.select(el)
        .attr("viewBox", [0.5, -60.5, width, height + 60]);

    let group = svg.append("g")
        .call(render, treemap(hierarchyData));

    function render(group, root) {
      const node = group
        .selectAll("g")
        .data(root.children.concat(root))
        .join("g");

      node.filter(d => d === root ? d.parent : d)
          .attr("cursor", "pointer")
          .on("click", (event, d) => {
            if (d === root) {
              dispatch('fileSelected', {
                file: d.parent.data
              });
              zoomout(root);
            }
            else {
              dispatch('fileSelected', {
                file: d.data
              });
              if (d.data.children && d.data.children.length > 0) {
                zoomin(d);
              }
            }
          });

      node.append("title")
          .text(d => `${name(d)}\n${format(d.data.lines)}`);

      node.append("rect")
          // .attr("id", d => d.leafUid = uuidv4())
          .attr("fill", d => d === root ? "#fff" : color(d.data.rating))
          .attr("stroke", "#fff")
          .on('mouseover', function() {
            tooltip.style('visibility', 'visible');
          })
          .on('mousemove', function(event, d) {
            tooltip
              .style('top', event.pageY - 10 + 'px')
              .style('left', event.pageX + 10 + 'px')
              .text(`Rating: ${d.data.rating}`);
          })
          .on('mouseout', function() {
            tooltip.style('visibility', 'hidden');
          });;

      node.append("clipPath")
        // .attr("id", d => d.clipUid = uuidv4())
        .append("use")
        // .attr("xlink:href", d => d.leafUid.href);

      node.append("text")
          .attr("clip-path", d => d.clipUid)
          .attr("font-weight", d => d === root ? "bold" : null)
          .selectAll("tspan")
          .data(d => [(d === root ? name(d) : d.data.name)].concat(format(d.data.lines)))
          .join("tspan")
          .attr("x", 3)
          .attr("y", (d, i, nodes) => `${(i === nodes.length - 1) * 0.3 + 1.1 + i * 0.9}em`)
          .attr("fill-opacity", (d, i, nodes) => i === nodes.length - 1 ? 0.7 : null)
          .attr("font-weight", (d, i, nodes) => i === nodes.length - 1 ? "normal" : null)
          .text(d => d);

      group.call(position, root);
    }

    function position(group, root) {
      group.selectAll("g")
          .attr("transform", d => d === root ? `translate(0,-60)` : `translate(${x(d.x0)},${y(d.y0)})`)
        .select("rect")
          .attr("width", d => d === root ? width : x(d.x1) - x(d.x0))
          .attr("height", d => d === root ? 60 : y(d.y1) - y(d.y0));
    }

    // When zooming in, draw the new nodes on top, and fade them in.
    function zoomin(d) {
      const group0 = group.attr("pointer-events", "none");
      const group1 = group = svg.append("g").call(render, d);

      x.domain([d.x0, d.x1]);
      y.domain([d.y0, d.y1]);

      svg.transition()
          .duration(750)
          .call(t => group0.transition(t).remove()
            .call(position, d.parent))
          .call(t => group1.transition(t)
            .attrTween("opacity", () => d3.interpolate(0, 1))
            .call(position, d));
    }

    // When zooming out, draw the old nodes on top, and fade them out.
    function zoomout(d) {
      const group0 = group.attr("pointer-events", "none");
      const group1 = group = svg.insert("g", "*").call(render, d.parent);

      x.domain([d.parent.x0, d.parent.x1]);
      y.domain([d.parent.y0, d.parent.y1]);

      svg.transition()
          .duration(750)
          .call(t => group0.transition(t).remove()
            .attrTween("opacity", () => d3.interpolate(1, 0))
            .call(position, d))
          .call(t => group1.transition(t)
            .call(position, d.parent));
    }

	});

  function flattenData(hierarchyData) {
    let keys = Object.keys(hierarchyData);
    if (hierarchyData.children === undefined) {
      return hierarchyData;
    }
    else {
      const children = [];
      Object.keys(hierarchyData.children).forEach((key) => {
        children.push(
          {
            ...flattenData(hierarchyData.children[key]),
            name: key
          }
        );
      });

      hierarchyData.children = children;
      return hierarchyData;
    }
  }

</script>


<svg bind:this={el} class="chart"></svg>

<style>
  .tooltip {
    position: relative;
    display: inline-block;
    /* border-bottom: 1px dotted black; */
  }
</style>