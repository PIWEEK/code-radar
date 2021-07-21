<svelte:head>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/css/uikit.min.css" />
  <!-- UIkit JS -->
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit-icons.min.js"></script>
</svelte:head>


<script lang="ts">
	import Treemap from './Treemap.svelte';

	async function getProjectData() {
		const res = await fetch(`http://localhost:8000/files`);
		const json = await res.json();

    if (res.ok) {
			return json;
		} else {
			throw new Error(json);
		}
	}
</script>

<main>
  {#await getProjectData()}
    <p>...waiting</p>
  {:then projectInfo}
    <h1 class="uk-heading-medium uk-heading-divider">
    <a href="{projectInfo.url}" target="_blank">{projectInfo.name}
    </a>
    </h1>

    <div class="chart">
      <Treemap data={projectInfo}/>
    </div>

<!--
    <div class="analytics">
      <ul>
        {#if selected.data.lines != undefined}
        <li>Lines: {selected.data.lines}</li>
        {/if}
        {#if selected.data.rating != undefined}
        <li>Rating: {selected.data.rating}</li>
        {/if}
      </ul>
    </div>
-->

<!--
    <ul class="uk-list uk-list-striped">
      {#each projectInfo.files as file, i}
        <li>
          {file.name}
        </li>
      {/each}
    </ul>
-->
  {:catch error}
    <p style="color: red">{error.message}</p>
  {/await}

</main>

<style>
	main {
		text-align: left;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}

	.chart :global(div) {
		font: 10px sans-serif;
		background-color: steelblue;
		text-align: right;
		padding: 3px;
		margin: 1px;
		color: white;
	}

	.chart {
    float: left;
		width: calc(80% + 2px);
		height: 400px;
		padding: 0;
		margin: 0 -1px 36px -1px;
		overflow: hidden;
	}

/* 
  .breadcrumbs {
		width: 100%;
		padding: 0.3rem 0.4rem;
		background-color: transparent;
		font-family: inherit;
		font-size: inherit;
		text-align: left;
		border: none;
		cursor: pointer;
		outline: none;
	}

	.breadcrumbs:disabled {
		cursor: default;
	}


  .analytics {
    width: calc(20% - 2px);
		height: 400px;
    float: right;
  }

	.node {
		position: absolute;
		width: 100%;
		height: 100%;
		background-color: white;
		overflow: hidden;
		pointer-events: all;
	}

	.node:not(.leaf) {
		cursor: pointer;
	}

	.contents {
		width: 100%;
		height: 100%;
		padding: 0.3rem 0.4rem;
		border: 1px solid white;
		background-color: hsl(240, 8%, 70%);
		color: white;
		border-radius: 4px;
		box-sizing: border-box;
	}

	.node:not(.leaf) .contents {
		background-color: hsl(240, 8%, 44%);
	}

	strong, span {
		display: block;
		font-size: 12px;
		white-space: nowrap;
		line-height: 1;
	} */

  /* .uk-list {
    clear: both;
  } */
</style>