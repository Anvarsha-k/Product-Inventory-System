<script>
  import Sidebar from "./components/Sidebar.svelte";
  import ProductForm from "./routes/ProductForm.svelte";
  import ProductList from "./routes/ProductList.svelte";
  import StockManage from "./routes/StockManage.svelte";
  import StockReport from "./routes/StockReport.svelte";

  let view = "create";
  function setView(v) { view = v; }
</script>

<style>
  .layout {
    display: flex;
  }

  .content {
    margin-left: 260px;
    padding: 20px;
    width: calc(100% - 260px);
  }
</style>

<div class="layout">
  <Sidebar {view} {setView} />

  <div class="content">
    {#if view === "create"}
      <ProductForm />
    {:else if view === "list"}
      <ProductList />
    {:else if view === "stock"}
      <StockManage />
    {:else}
      <StockReport />
    {/if}
  </div>
</div>
