<script>
  export let view;
  export let setView;
</script>

<style>
  .sidebar {
    width: 230px;
    background: #1e1e2f;
    height: 100vh;
    color: #fff;
    padding: 20px;
    position: fixed;
    top: 0;
    left: 0;
  }

  h2 {
    margin: 0 0 25px 0;
    font-size: 20px;
    font-weight: 600;
  }

  .nav-item {
    padding: 10px 0;
    cursor: pointer;
    font-size: 15px;
  }

  .nav-item:hover,
  .active {
    color: #00b7ff;
  }
</style>

<div class="sidebar">
  <h2>Inventory</h2>

  <div class="nav-item" class:active={view === "create"} on:click={() => setView("create")}>Create Product</div>
  <div class="nav-item" class:active={view === "list"} on:click={() => setView("list")}>Product List</div>
  <div class="nav-item" class:active={view === "stock"} on:click={() => setView("stock")}>Stock Manage</div>
  <div class="nav-item" class:active={view === "report"} on:click={() => setView("report")}>Stock Report</div>
</div>
