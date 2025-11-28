<script>
  let product = {
    product_id: "",
    product_code: "",
    product_name: "",
    product_image: "",
    variants: [],
    sub_variants: []
  };

  function addVariant() {
    product.variants.push({ name: "", options: [""] });
    product = { ...product };
  }

  function addOption(v) {
    v.options.push("");
    product = { ...product }; 
  }

  function addSubVariant() {
    product.sub_variants.push({ option_ids: [], sku: "", stock: "0" });
    product = { ...product };
  }

  async function submit() {
    try {
      product.product_id = Number(product.product_id);

      const res = await fetch("http://localhost:8080/api/products", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(product)
      });

      const json = await res.json();

      if (!res.ok) {
        alert("Error: " + JSON.stringify(json));
        return;
      }

      alert("Product Created!");

      product = {
        product_id: "",
        product_code: "",
        product_name: "",
        product_image: "",
        variants: [],
        sub_variants: []
      };
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
    margin-bottom: 20px;
    box-shadow: 0 3px 10px rgba(0,0,0,0.04);
  }

  input {
    width: 100%;
    padding: 10px;
    margin-top: 6px;
    margin-bottom: 12px;
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
    font-size: 14px;
    margin-top: 8px;
  }

  button.secondary {
    background: #444;
  }

  .block {
    padding: 15px;
    background: #f7f7f7;
    border-radius: 8px;
    margin-top: 12px;
  }

  h3, h4 {
    margin-bottom: 10px;
    margin-top: 0;
  }

</style>

<div class="card">
  <h2>Create Product</h2>


  <h3>Basic Information</h3>
  <input placeholder="Product ID" bind:value={product.product_id} />
  <input placeholder="Product Code" bind:value={product.product_code} />
  <input placeholder="Product Name" bind:value={product.product_name} />
  <input placeholder="Product Image URL" bind:value={product.product_image} />
</div>

<!-- VARIANTS -->
<div class="card">
  <h3>Variants</h3>

  <button on:click={addVariant}>+ Add Variant</button>

  {#each product.variants as v, i}
    <div class="block">
      <h4>Variant {i + 1}</h4>

      <input placeholder="Variant Name (e.g., Size, Color)" bind:value={v.name} />

      <h4>Options</h4>
      {#each v.options as op, j}
        <input placeholder="Option Value (e.g., S, M, Red)" bind:value={v.options[j]} />
      {/each}

      <button class="secondary" on:click={() => addOption(v)}>+ Add Option</button>
    </div>
  {/each}
</div>

<!-- SUBVARIANTS -->
<div class="card">
  <h3>Sub Variants (SKU Combinations)</h3>

  <button on:click={addSubVariant}>+ Add Subvariant</button>

  {#each product.sub_variants as sv, i}
    <div class="block">
      <h4>Subvariant {i + 1}</h4>
      <input placeholder="SKU" bind:value={sv.sku} />
      <input placeholder="Stock" bind:value={sv.stock} />

      <input
        placeholder="Option IDs (comma separated)"
        value={sv.option_ids.join(", ")}
        on:input={(e) =>
          sv.option_ids = e.target.value.split(",").map(s => s.trim())
        }
      />
    </div>
  {/each}
</div>

<button on:click={submit}>Create Product</button>
