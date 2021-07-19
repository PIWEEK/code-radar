<svelte:head>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/css/uikit.min.css" />
  <!-- UIkit JS -->
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/uikit@3.7.1/dist/js/uikit-icons.min.js"></script>
</svelte:head>


<script lang="ts">
	// export let name: string;
	async function getProjectData() {
		const res = await fetch(`https://mocki.io/v1/52e7b578-b952-47ca-8c58-618aee7bf41b`);
		const text = await res.json();

		if (res.ok) {
			return text;
		} else {
			throw new Error(text);
		}
	}

  function buildPath(path, name, extension) {
    let fullPath = "";
    if(path.length >= 0) {
      fullPath = path.join(`/`);
    }
    return `${fullPath}${name}.${extension}`;
  }
</script>

<main>
  {#await getProjectData()}
    <p>...waiting</p>
  {:then projectInfo}
    <h1 class="uk-heading-medium uk-heading-divider">{projectInfo.name} - {projectInfo.url}</h1>
    <ul class="uk-list uk-list-striped">
      {#each projectInfo.files as file, i}
        <li>
          {buildPath(file.path, file.name, file.extension)}
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
</style>