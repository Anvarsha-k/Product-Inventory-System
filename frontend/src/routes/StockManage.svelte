<script>
  import { onMount } from "svelte";

  let subVariants = [];
  let selected = "";
  let qty = "1";
  let type = "IN";

  async function loadSubs() {
    try {
      const res = await fetch("http://localhost:8080/api/Listproducts?page=1&limit=200");
      const data = await res.json();

      if (!res.ok) {
        alert("Error: " + JSON.stringify(data));
        return;
      }

      // IMPORTANT FIX: handle products with no sub_variants
      subVariants = data.data.products.flatMap(p =>
        (p.sub_variants || []).map(sv => ({
          id: sv.id,
          sku: sv.sku,
          productName: p.product_name
        }))
      );

    } catch (err) {
      alert("Network error: " + err.message);
    }
  }

  onMount(loadSubs);

  async function adjustStock() {
    if (!selected) {
      alert("Please select a Subvariant");
      return;
    }

    try {
      const body = {
        sub_variant_id: selected,
        quantity: qty,
        transaction_type: type
      };

      const res = await fetch("http://localhost:8080/api/stock/adjust", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
      });

      const json = await res.json();

      if (!res.ok) {
        alert("Error: " + JSON.stringify(json));
        return;
      }

      alert("Stock Updated Successfully!");
      qty = "1";

    } catch (err) {
      alert("Network error: " + err.message);
    }
  }
</script>

<style>
  .card {
    background: #fff;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 3px 10px rgba(0,0,0,0.05);
    max-width: 500px;
  }

  select, input {
    width: 100%;
    padding: 10px;
    margin-top: 6px;
    margin-bottom: 15px;
    border-radius: 6px;
    border: 1px solid #d0d0d0;
    font-size: 14px;
  }

  button {
    padding: 10px 16px;
    background: #007bff;
    color: #fff;
    border-radius: 6px;
    cursor: pointer;
    border: none;
    width: 100%;
    font-size: 15px;
  }
</style>

<div class="card">
  <h2>Stock Manage</h2>

  <label>Select Subvariant (SKU)</label>
  <select bind:value={selected}>
    <option value="">-- Select Subvariant --</option>
    {#each subVariants as sv}
      <option value={sv.id}>{sv.productName} - {sv.sku}</option>
    {/each}
  </select>

  <label>Quantity</label>
  <input type="number" bind:value={qty} min="1" />

  <label>Transaction Type</label>
  <select bind:value={type}>
    <option value="IN">IN (Add Stock)</option>
    <option value="OUT">OUT (Remove Stock)</option>
  </select>

  <button on:click={adjustStock}>Update Stock</button>
</div>
