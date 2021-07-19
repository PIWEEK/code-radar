<svelte:head>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/css/uikit.min.css" />
  <!-- UIkit JS -->
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit-icons.min.js"></script>
</svelte:head>


<script lang="ts">
	// export let name: string;

	import * as Pancake from '@sveltejs/pancake';
	import * as d3 from 'd3-hierarchy';
	import { tweened } from 'svelte/motion';
	import * as eases from 'svelte/easing';
	import { fade } from 'svelte/transition';
	import * as yootils from 'yootils';
	import Treemap from './Treemap.svelte';
	import data from './data.js';

  const extents = tweened(undefined, {
    easing: eases.cubicOut,
    duration: 600
  });

  let selected = undefined;
  let root = undefined;

  $: $extents = {
    x1: selected ? selected.x0 : 0,
    x2: selected ? selected.x1 : 0,
    y1: selected ? selected.y1 : 0,
    y2: selected ? selected.y0 : 0
  };


  const breadcrumbs = node => {
    const crumbs = [];
    while (node) {
      crumbs.unshift(node.data.name)
      node = node.parent;
    }

    return crumbs.join('/');
  };

  const is_visible = (a, b) => {
    console.log("is_visible", a ,b)
    while (b) {
      if (a.parent === b) return true;
      b = b.parent;
    }

    return false;
  };

  const select = node => {
    console.log("select", 1)
    while (node.parent && node.parent !== selected) {
      console.log("select", 2)
      node = node.parent;
    }
    console.log("select", 3)

    if (node && node.children) selected = node;
    console.log("select", selected)
  };

	async function getProjectData() {
		const res = await fetch(`http://localhost:8000/files?path=`);
		const text = await res.json();

		if (res.ok) {
      const treemap = d3.treemap();

      // const hierarchy = d3.hierarchy(data)
      //   .sum(d => d.value)
      //   .sort((a, b) => b.value - a.value)

      const hierarchyData = {
        name: "/",
        children: text.files
      }

      const hierarchy = d3.hierarchy(hierarchyData)
        .sum(d => d.Lines)
        .sort((a, b) => b.Lines - a.Lines)

      console.log("hierarchy", hierarchy)

      root = treemap(hierarchy);
      
      selected = root;

			return text;
		} else {
			throw new Error(text);
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

    <button class="breadcrumbs" disabled="{!selected.parent}" on:click="{() => selected = selected.parent}">
      {breadcrumbs(selected)}
    </button>

    <div class="chart">
      <Pancake.Chart x1={$extents.x1} x2={$extents.x2} y1={$extents.y1} y2={$extents.y2}>
        <Treemap {root} let:node>
          {#if is_visible(node, selected)}
            <div
              transition:fade={{duration:400}}
              class="node"
              class:leaf={!node.children}
              on:click="{() => select(node)}"
            >
              <div class="contents">
                <strong>{node.data.name}</strong>
                <span>{yootils.commas(node.value)}</span>
              </div>
            </div>
          {/if}
        </Treemap>
      </Pancake.Chart>
    </div>

    <ul class="uk-list uk-list-striped">
      {#each projectInfo.files as file, i}
        <li>
          {file.name}
        </li>
      {/each}
    </ul>
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

	/* h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	} */

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}

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

	.chart {
		width: calc(100% + 2px);
		height: 400px;
		padding: 0;
		margin: 0 -1px 36px -1px;
		overflow: hidden;
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
	}
</style>