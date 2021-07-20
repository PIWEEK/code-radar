<svelte:head>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/css/uikit.min.css" />
  <!-- UIkit JS -->
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit-icons.min.js"></script>
</svelte:head>


<script lang="ts">
	// export let name: string;
  let dimension = 'lines';

	import * as Pancake from '@sveltejs/pancake';
	import * as d3 from 'd3-hierarchy';
	import { tweened } from 'svelte/motion';
	import * as eases from 'svelte/easing';
	import { fade } from 'svelte/transition';
	import * as yootils from 'yootils';
	import Treemap from './Treemap.svelte';
	import data from './data.js';

  const treemap = d3.treemap();

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

  const hierarchyData = {
    name: "",
    children: {}
  };

  function breadcrumbs(node) {
    const crumbs = [];
    while (node) {
      crumbs.unshift(node.data.name)
      node = node.parent;
    }
		const b = crumbs.join('/');
    return b.startsWith('/') ? b.substring(1) : b;
  };

  function is_visible(a, b) {
    while (b) {
      if (a.parent === b) return true;
      b = b.parent;
    }

    return false;
  };

  async function select(node) {
    while (node.parent && node.parent !== selected) {
      node = node.parent;
    }
    if (node && node.children) selected = node;
    selected = node;
  };

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

	async function getProjectData() {
		const res = await fetch(`http://localhost:8000/files`);
		const json = await res.json();

    const maxLines = json.files.reduce((prev, current) => (prev.lines > current.lines) ? prev : current).lines;
    const minLines = maxLines * .10;
    const maxRating = json.files.reduce((prev, current) => (prev.rating > current.rating) ? prev : current).rating;
    const minRating = maxRating * .10;

		if (res.ok) {
      json.files.forEach((f) => {
        let path = f.directory  ? f.directory.split("/") : [];
        if (f.isDirectory) {
          path = path.concat([f.name])
        }

        let d = hierarchyData;
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

        if (!f.isDirectory) {
          d.children[f.name] = {
            lines: f.lines,
            rating: f.rating
          };
        } else{
          d.lines = f.lines;
          d.rating = f.rating;
        }
      });

      flattenData(hierarchyData);

      const hierarchy = d3.hierarchy(hierarchyData)
        .sum(d => d.lines)
        .sort((a, b) => b.lines - a.lines)

      root = treemap(hierarchy);
      selected = root;

			return json;
		} else {
			throw new Error(json);
		}
	}

  function onChangeDimension(event) {
		const val = event.currentTarget.value;
    const hierarchy = d3.hierarchy(hierarchyData)
      .sum(d => d[val])
      .sort((a, b) => b[val] - a[val])

    root = treemap(hierarchy);
    selected = root;

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

    <label>
      <input type=radio on:change={onChangeDimension} bind:group={dimension} name="dimension" value={'lines'}>
      lines
    </label>
    <label>
      <input type=radio on:change={onChangeDimension} bind:group={dimension} name="dimension" value={'rating'}>
      rating
    </label>

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
    float: left;
		width: calc(80% + 2px);
		height: 400px;
		padding: 0;
		margin: 0 -1px 36px -1px;
		overflow: hidden;
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
	}

  /* .uk-list {
    clear: both;
  } */
</style>