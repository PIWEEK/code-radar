<script>
 import { onMount, afterUpdate } from 'svelte';
 import * as d3 from 'd3';

 export let file;
 export let firstCommit;
 export let lastCommit;

 function createActivityChart(from, to, data) {
   // set the dimensions and margins of the graph
   const margin = {top: 10, right: 30, bottom: 70, left: 40};
   const width = 460 - margin.left - margin.right;
   const height = 250 - margin.top - margin.bottom;
   const viewbox = [0, 0,
                    width + margin.left + margin.right,
                    height + margin.top + margin.bottom].join(" ")

   // append the svg object to the body of the page
   const svg =
     d3.select(".activity-chart")
       .append("div")
       .classed("graph", true) 
       .append("svg")
       .attr("preserveAspectRatio", "xMinYMin meet")
       .attr("viewBox", viewbox)
       .classed("graph-content", true)
       .append("g")
       .attr("transform", `translate(${margin.left},${margin.top})`);

   const x =
     d3.scaleTime()
       .domain([from, to])
       .rangeRound([0, width])
       .nice();

   svg.append("g")
      .attr("transform", `translate(0, ${height})`)
      .call(d3.axisBottom(x)
        //.tickFormat(d3.timeFormat("%Y-%m-%d"))
      )
      .selectAll("text")	
      .style("text-anchor", "end")
      .attr("dx", "-.8em")
      .attr("dy", ".15em")
      .attr("transform", "rotate(-65)")
   ;

   const maxValue = data.reduce((acc,d) => Math.max(acc,d.Added + d.Deleted), -Infinity)
   const y = d3.scaleLinear()
               .domain([0, maxValue + 10])
               .range([height, 0]);

   svg.append("g")
      .call(d3.axisLeft(y));


   svg.selectAll("barchart")
      .data(data)
      .enter()
      .append("rect")
      .attr("x", function(d) { return x(Date.parse(d.Date)); })
      .attr("y", function(d) { return y(d.Added + d.Deleted); })
      .attr("width", 2)
      .attr("height", function(d) { return height - y(d.Added + d.Deleted); })
      .attr("fill", "#69b3a2")
 }

 function updateActivityChart(data) {
   document.querySelectorAll(".activity-chart *").forEach((n)=>n.remove())
   createActivityChart(firstCommit, lastCommit, file.history);
 }

 onMount(async () => {
   createActivityChart(firstCommit, lastCommit, file.history);
 });
 
 afterUpdate(async () => {
   updateActivityChart(firstCommit, lastCommit, file.history);
 });

</script>

<div class="activity-chart">
</div>

<style>
 .graph {
   display: inline-block;
   position: relative;
   width: 100%;
   padding-bottom: 100%; /* aspect ratio */
   vertical-align: top;
   overflow: hidden;
 }

 .graph-content {
   display: inline-block;
   position: absolute;
   top: 10px;
   left: 0;
 }
</style>
