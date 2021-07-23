<script>
 import ActivityChart from "./ActivityChart.svelte";
 import ActivityList from "./ActivityList.svelte";
 import OwnersChart from "./OwnersChart.svelte";

 import { onMount, afterUpdate } from 'svelte';
 import * as d3 from 'd3';

 export let project;
 export let file;
 export let userColors;
 export let firstCommit;
 export let lastCommit;
</script>

<div class="detail">

  <div class="header">
    <div class="header-title">
      <img class="logo" src="/favicon.svg"/>
      {#if file.name == "."}
        <div><a href={project.url}>{project.name}</a></div>
      {:else if !file.directory || file.directory === "."}
        <div><a href={project.url}>{project.name}</a> / {file.name}</div>
      {:else}
        <div><a href={project.url}>{project.name}</a> / {file.directory} / {file.name}</div>
      {/if}
    </div>
    <div class="header-detail">{file.lines} lines. {file.history.length} changes. Rating: {Number(file.rating * 100).toFixed(2)}%</div>
  </div>

  <div class="content">
    <h3>Owners</h3>
    <OwnersChart file={file} userColors={userColors}/>

    <h3>Activity</h3>
    <ActivityChart file={file} firstCommit={firstCommit} lastCommit={lastCommit}/>

    <h3>Changes</h3>
    <ActivityList file={file}/>
  </div>
</div>

<style>
 .detail {
   display: flex;
   flex-direction: column;
   height: 100%;
   margin: 0rem;
   overflow: hidden;
   padding: 0;
 }

 h3 {
   color: #1e87f0;
   margin: 0;
   border-bottom: 1px solid #1e87f0;
   margin-bottom: 1.5rem;
   font-size: 1rem;
   font-variant: all-petite-caps;
   font-weight: 600;
 }

 .header {
   padding: 1rem;
 }

 .content {
   flex: 1;
   padding: 1rem;
   overflow-y: scroll;
 }
 
 .header a {
   font-size: 1.5rem;
   font-weight: 600;
 }

 .header-title {
   display: flex;
   align-items: center;
 }

 .header-title .logo {
   width: 1.5rem;
   margin-right: 0.5rem;
 }
 
 .header-detail {
   font-size: 0.8rem;
   font-weight: 300;
 }
 

</style>

