<script>
 import { onMount, afterUpdate } from 'svelte';
 import * as d3 from 'd3';

 export let file;
 export let userColors;

 function userColor(user) {
   if (user === "other") {
     return "#666";
   }
   return userColors[user];
 }
 
 function createOwnersChart(data) {
   const width = 100;
   const height = 100;

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

   
   const pieData = pie(processData(data))

   svg.selectAll('piechart')
      .data(pieData)
      .join('path')
      .attr('d', d3.arc().innerRadius(0).outerRadius(radius))
      .attr('fill', (d) => userColor(d.data[0]))
      .attr("stroke", "none");     
 }

 function processData(data) {
   let others = 0.0;
   let result = [];

   // Puts together all the "minor" percents into one single category "other"
   Object.entries(data)
         .forEach(([user, percent]) => {
           if (percent > .01) {
             result.push([user, percent]);
           } else {
             others += percent;
           }
         });

   if (others > .01) {
     result.push(["other", others]);
   }

   return result;
 }

 function updateOwnersChart(data) {
   document.querySelectorAll(".owners-chart *").forEach((n)=>n.remove());
   createOwnersChart(file.owners);
 }

 function sortedOwners(owners) {
   const entries = processData(owners)
   return entries.map(([user, value]) => ({
     user: user,
     percent: Number(value * 100).toFixed(2),
     color: userColor(user)
   })).sort((a, b) => b.percent - a.percent);
 }

 onMount(async () => {
   createOwnersChart(file.owners, userColors);
 });
 
 afterUpdate(async () => {
   updateOwnersChart(file.owners, userColors);
 });
</script>


<div class="owner-info">
  <ul class="user-list">
    {#each sortedOwners(file.owners) as owner}
      <li class="user-list-entry" title="{owner.user} ({owner.percent}%)">
        <div class="user-list-color" style="background-color: {owner.color}"></div>
        <div class="user-list-data">
          <span class="user-list-name">{owner.user}</span> <span class="user-list-percent">({owner.percent}%)<span>
        </div>
      </li>
    {/each}
  </ul>

  <!-- Placeholder for the Chart -->
  <div class="owners-chart">
  </div>
</div>

<style>
 .user-list {
   font-size: 70%;
   list-style: none;
   padding: 0;
   min-width: 15.5rem;
   display: flex;
   flex-direction: column;
   justify-content: center;
   padding: 8px;
 }

 .user-list-entry {
   display: flex;
   align-items: center;
 }

 .user-list-data {
   white-space: nowrap;
   overflow: hidden;
   text-overflow: ellipsis;
 }
 .user-list-color {
   width: 8px;
   height: 8px;
   margin-right: 8px
 }

 .owner-info {
   display: flex;
   align-items: center;
 }

 .owners-chart {
   width: 100%;
   height: 100%;
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
