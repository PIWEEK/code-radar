<script>
 import { onMount, afterUpdate } from 'svelte';
 import * as d3 from 'd3';

 export let file;
 export let userColors;
 export let firstCommit;
 export let lastCommit;

 function createOwnersChart(data, dataColor) {
   const width = 100;
   const height = 50;

   const radius = Math.min(width, height) / 2;

   const svg =
     d3.select(".owners-chart")
       .append("div")
       .classed("graph", true) 
       .append("svg")
       .attr("preserveAspectRatio", "xMinYMin meet")
       .attr("viewBox", [0, 0, width, height].join(" "))
       .classed("graph-content", true)
       .append("g")
       .attr("transform", `translate(${width/2}, ${height/2})`);

   const pie = d3.pie().value((d) => d[1])
   const pieData = pie(Object.entries(data))

   svg
     .selectAll('piechart')
     .data(pieData)
     .join('path')
     .attr('d', d3.arc().innerRadius(0).outerRadius(radius))
     .attr('fill', (d) => dataColor[d.data[0]])
     .attr("stroke", "none");     
 }

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

 function updateOwnersChart(data) {
   document.querySelectorAll(".owners-chart *").forEach((n)=>n.remove());
   createOwnersChart(file.owners, userColors);
 }

 function updateActivityChart(data) {
   document.querySelectorAll(".activity-chart *").forEach((n)=>n.remove())
   createActivityChart(firstCommit, lastCommit, file.history);
 }


 onMount(async () => {
   createOwnersChart(file.owners, userColors);
   createActivityChart(firstCommit, lastCommit, file.history);
 });
 
 afterUpdate(async () => {
   updateOwnersChart(file.owners, userColors);
   updateActivityChart(firstCommit, lastCommit, file.history);
 });

</script>

<div class="detail">
  <div>Name: {file.name}</div>
  <div>Lines: {file.lines}</div>
  <div>Rating: {file.rating}</div>

  <div class="owners-chart">
  </div>

  <div class="activity-chart">
  </div>
</div>

<style>
 .detail {
   border: 1px solid black;
   margin: 0.5rem;
   padding: 0.5rem;
 }
 .owners-chart {
   
 }

 .activity-chart {
   
 }

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

