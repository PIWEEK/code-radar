<script>
 export let file;
</script>

<ul class="file-history uk-list uk-list-striped">
  {#each file.history as history, i}
    <li>
      [{history.Date}] {history.User} +{history.Added} -{history.Deleted}
    </li>
  {/each}
</ul>

<style>
 .file-history {
   font-size: 0.8rem;
 }
</style>
