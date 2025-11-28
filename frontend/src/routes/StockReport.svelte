<script>
  let from = "";
  let to = "";
  let report = null;

  async function fetchReport() {
    try {
      const url = new URL("http://localhost:8080/api/stock/report");

      if (from) url.searchParams.append("from", new Date(from).toISOString());
      if (to)   url.searchParams.append("to", new Date(to).toISOString());

      const res = await fetch(url);
      const json = await res.json();

      if (!res.ok) {
        alert("Error: " + JSON.stringify(json));
        return;
      }

      report = json.data;

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
    box-shadow: 0 3px 10px rgba(0,0,0,0.05);
  }

  input {
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
    font-size: 15px;
    width: 100%;
  }

  .summary-box {
    background: #f7f7f7;
    padding: 14px;
    margin: 10px 0;
    border-radius: 8px;
    border-left: 4px solid #007bff;
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
    background: #f5f5f5;
    text-align: left;
  }
</style>

<div class="card">
  <h2>Stock Report</h2>

  <label>From Date</label>
  <input type="datetime-local" bind:value={from} />

  <label>To Date</label>
  <input type="datetime-local" bind:value={to} />

  <button on:click={fetchReport}>Get Report</button>
</div>

{#if report}
  <div class="card">
    <h3>Summary</h3>

    {#each report.summary as s}
      <div class="summary-box">
        <strong>Product:</strong> {s.product_id} <br />
        <strong>Stock In:</strong> {s.stock_in} <br />
        <strong>Stock Out:</strong> {s.stock_out} <br />
        <strong>Net:</strong> {s.net}
      </div>
    {/each}
  </div>

  <div class="card">
    <h3>Transactions</h3>

    <table>
      <thead>
        <tr>
          <th>Date</th>
          <th>Type</th>
          <th>Quantity</th>
          <th>Subvariant</th>
        </tr>
      </thead>

      <tbody>
        {#each report.transactions as t}
          <tr>
            <td>{t.transaction_date}</td>
            <td>{t.transaction_type}</td>
            <td>{t.quantity}</td>
            <td>{t.sub_variant_id}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}
