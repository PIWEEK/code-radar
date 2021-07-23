<svelte:head>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/css/uikit.min.css" />
  <!-- UIkit JS -->
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit-icons.min.js"></script>
</svelte:head>


<script lang="ts">
  import Treemap from './Treemap.svelte';
  import Detail from './detail/Detail.svelte';

  let selected = undefined;
  let innerHeight = undefined;
  let innerWidth = undefined;

  let userColors;
  let firstCommit;
  let lastCommit;

	async function getProjectData() {
		const res = await fetch(`/files`);
		const json = await res.json();

    if (res.ok) {

      userColors = calculateUserColors(json);
      firstCommit = Date.parse(json["firstCommit"]);
      lastCommit = Date.parse(json["lastCommit"]);

      // Select root by default
      selected = json["files"].find((it)=>it.name === ".");
			return json;
		} else {
			throw new Error(json);
		}
	}

  function calculateUserColors(json) {
    const users = json["users"]

    const colors = [
      "#3366cc", "#dc3912", "#ff9900", "#109618", "#990099",
      "#0099c6", "#dd4477", "#66aa00", "#b82e2e", "#316395",
      "#994499", "#22aa99", "#aaaa11", "#6633cc", "#e67300",
      "#8b0707", "#651067", "#329262", "#5574a6", "#3b3eac",
    ];

    var userColors =
      users.map((it, i) => [it, colors[i % colors.length] ])
           .reduce((acc, [k, v]) => ({...acc, [k]: v}), {});

    return userColors;
  }

 function handleFileSelected(file) {

    selected = file.detail.file;
  }
</script>

<svelte:window bind:innerHeight={innerHeight} bind:innerWidth={innerWidth}/>

<main>
  {#await getProjectData()}
    <p>...waiting</p>
  {:then projectInfo}
    <div class="detail">
      {#if selected}
        <Detail project={projectInfo}
                file={selected}
                userColors={userColors}
                firstCommit={firstCommit}
                lastCommit={lastCommit} />
      {/if}
    </div>

    <div class="chart">
      <Treemap data={projectInfo} width={innerWidth - 544} height={innerHeight - 40} on:fileSelected={handleFileSelected}/>
    </div>

  {:catch error}
    <p style="color: red">{error.message}</p>
  {/await}

</main>

<style>
 main {
   display: flex;
   flex-direction: row;
   width: 100%;
   height: 100%;
 }

 @media (min-width: 640px) {
	 main {
		 max-width: none;
	 }
 }

 .detail {
   width: 34rem;
   box-shadow: 1px 0px 6px 2px #00000030;
 }

 .chart {
   flex: 1;
   padding: 1rem 2rem;
 }

</style>
