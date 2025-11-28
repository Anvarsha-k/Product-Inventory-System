<script>
  import { onMount } from "svelte";

  let products = [];
  let page = 1;
  let limit = 10;
  let total = 0;

  async function loadProducts() {
    try {
      const res = await fetch(`http://localhost:8080/api/Listproducts?page=${page}&limit=${limit}`);
      const data = await res.json();

      if (!res.ok) {
        alert("Error: " + JSON.stringify(data));
        return;
      }

      products = data.data.products || [];
      total = data.data.total || 0;

    } catch (err) {
      alert("Network error: " + err.message);
    }
  }

  onMount(loadProducts);
</script>

<style>
  .card {
    background: #fff;
    padding: 20px;
    border-radius: 10px;
    margin-bottom: 20px;
    box-shadow: 0 3px 10px rgba(0,0,0,0.05);
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 15px;
  }

  th, td {
    padding: 12px;
    border-bottom: 1px solid #ececec;
    font-size: 14px;
  }

  th {
    text-align: left;
    background: #f5f5f5;
  }

  .badge {
    display: inline-block;
    background: #007bff;
    padding: 4px 8px;
    color: #fff;
    border-radius: 4px;
    font-size: 12px;
    margin-right: 4px;
  }

  button {
    padding: 8px 14px;
    background: #007bff;
    color: #fff;
    border-radius: 6px;
    border: none;
    cursor: pointer;
    margin: 10px 5px 0 0;
  }
</style>

<div class="card">
  <h2>Product List</h2>

  {#if products.length === 0}
    <p>No products found.</p>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Code</th>
          <th>Total Stock</th>
          <th>Variants</th>
          <th>Subvariants</th>
        </tr>
      </thead>

      <tbody>
        {#each products as p}
          <tr>
            <td>{p.product_name}</td>
            <td>{p.product_code}</td>
            <td>{p.total_stock}</td>

            <td>
              {#if p.variants}
                {#each p.variants as v}
                  <div>
                    <strong>{v.name}:</strong>
                    {#if v.options}
                      {v.options.map(o => o.value).join(", ")}
                    {/if}
                  </div>
                {/each}
              {/if}
            </td>

            <td>
              {#if p.sub_variants}
                {#each p.sub_variants as sv}
                  <div class="badge">{sv.sku} ({sv.stock})</div>
                {/each}
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}

  <button on:click={() => { if (page > 1) { page--; loadProducts(); }}}>
    Previous
  </button>

  <button on:click={() => { page++; loadProducts(); }}>
    Next
  </button>

  <p>Page: {page}</p>
</div>
