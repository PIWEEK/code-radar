<script>
 export let file;
</script>

<ul>
  {#if file?.name != undefined}
    <li>Name: {file?.name}</li>
  {/if}      
  {#if file?.lines != undefined}
    <li>Lines: {file?.lines}</li>
  {/if}
  {#if file?.rating != undefined}
    <li>Rating: {file?.rating}</li>
  {/if}
</ul>

<style>
</style>

