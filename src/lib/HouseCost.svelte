<script lang="ts">
    import { onMount } from "svelte";

    // Define types for the house cost breakdown
    type CostItem = {
        name: string;
        category: string;
        value: number;
        estimate: boolean;
        info: string;
    };

    // Reactive variables (Svelte reactivity)
    let house_price: number = 300000;
    let is_fisrt_house: boolean = true;
    let is_sold_by_builder: boolean = true;

    //20.5 is the number of mq to have a unit of "vani" for A/2 A/3 class
    //18.89 is a factor to obtain the estimate teriffe catastali which is divided by €/mq
    let rendita_catastale: number = house_price / 18.89 / 20.5;
    let valore_catastale: number = 0;

    let vat: number = 4;
    let vatAmount: number = (vat * house_price) / 100;
    let preliminare_pages: number = 4;
    let registrazione_preliminare: number = 200 + (16 * preliminare_pages) / 4;
    let imposta_di_registro: number = 200;
    let imposta_ipotecaria: number = 200;
    let imposta_catastale: number = 200;

    let agencyFee: number = 3;
    let agencyAmount: number = (agencyFee * house_price) / 100;
    let notaryAmount: number = house_price * 0.00085 * 1.22;

    let imposta_di_bollo = 230;
    let tassa_ipotecaria = 90;
    let tassa_archivio = 60;
    let visure_ipotecarie = 180;

    // Derived calculations (automatically update when dependencies change)
    //estimate to calculate the tassa_archivio cost, min 27.1 max 139.4

    $: tassa_archivio =
        house_price * 0.0001614 < 27.1
            ? 27.1
            : house_price * 0.0001614 > 139.4
              ? 139.4
              : house_price * 0.0001614;

    $: visure_ipotecarie = house_price * 0.000489 * 1.22; //estimate to calculate the visure_ipotecarie cost, can vary between notary

    $: rendita_catastale = house_price / 18.89 / 20.5;
    $: valore_catastale = is_sold_by_builder
        ? 0
        : is_fisrt_house
          ? rendita_catastale * 1.05 * 110
          : rendita_catastale * 1.05 * 120;

    $: imposta_di_registro = is_fisrt_house
        ? is_sold_by_builder
            ? 200
            : (2 * valore_catastale) / 100 > 1000
              ? (2 * valore_catastale) / 100
              : 1000
        : is_sold_by_builder
          ? 200
          : (9 * valore_catastale) / 100 > 1000
            ? (9 * valore_catastale) / 100
            : 1000;

    $: imposta_ipotecaria = is_fisrt_house ? 200 : 50;
    $: imposta_catastale = is_fisrt_house ? 200 : 50;
    $: vat = is_fisrt_house
        ? is_sold_by_builder
            ? 4
            : 0
        : is_sold_by_builder
          ? 10
          : 0;

    $: vatAmount = (vat * house_price) / 100;
    $: notaryAmount = house_price * 0.0042 * 1.22;

    $: agencyAmount = ((agencyFee * house_price) / 100) * 1.22;

    $: total_cost =
        house_price +
        agencyAmount +
        vatAmount +
        imposta_di_registro +
        imposta_ipotecaria +
        imposta_catastale +
        notaryAmount +
        tassa_archivio +
        tassa_ipotecaria +
        visure_ipotecarie +
        imposta_di_bollo +
        registrazione_preliminare;

    // Chart Data
    let chartData: CostItem[] = [];

    onMount(() => {
        updateChart();
    });

    function updateChart(): void {
        chartData = [
            {
                name: "Property Price",
                category: "House price",
                value: Math.round(house_price),
                estimate: false,
                info: "House buy price",
            },
            {
                name: "Agency",
                category: "Agency",
                value: Math.round(agencyAmount),
                estimate: false,
                info: "Agency fee",
            },
            {
                name: "VAT",
                category: "Tax",
                value: Math.round(vatAmount),
                estimate: false,
                info: "IVA dovuta al venditore, incide solamente se il venditore è passibile di IVA.",
            },
            {
                name: "Imposta di registro",
                category: "Tax",
                value: Math.round(imposta_di_registro),
                estimate: false,
                info: "Imposta di registro, tassa che varia in base al tipo di venditore.",
            },
            {
                name: "Imposta ipotecaria",
                category: "Tax",
                value: Math.round(imposta_ipotecaria),
                estimate: false,
                info: "Imposta ipotecaria, tassa fissa che varia in base al tipo di venditore.",
            },
            {
                name: "Imposta catastale",
                category: "Tax",
                value: Math.round(imposta_catastale),
                estimate: false,
                info: "Imposta catastale, tassa fissa che varia in base al tipo di venditore.",
            },
            {
                name: "Registrazione preliminare",
                category: "Tax",
                value: Math.round(registrazione_preliminare),
                estimate: false,
                info: "Registrazione preliminare, spesa da sostenere durante il preliminare",
            },
            {
                name: "Notary",
                category: "Notary",
                value: Math.round(notaryAmount),
                estimate: true,
                info: "Compenso onorario notaio",
            },
            {
                name: "Tassa di archivio",
                category: "Notary",
                value: Math.round(tassa_archivio),
                estimate: true,
                info: "Tassa di archivio, varia in relazione all'importo dell'immobile.",
            },
            {
                name: "Tassa ipotecaria",
                category: "Notary",
                value: Math.round(tassa_ipotecaria),
                estimate: false,
                info: "Tassa ipotecaria, tassa fissa durante il rogito.",
            },
            {
                name: "Visure ipotecarie",
                category: "Notary",
                value: Math.round(visure_ipotecarie),
                estimate: true,
                info: "Visure ipotecarie, spese variabile da sostenere per la consultazione delle visure ipotecarie da parte del notaio durante il rogito.",
            },
            {
                name: "Imposta di bollo",
                category: "Notary",
                value: Math.round(imposta_di_bollo),
                estimate: false,
                info: "Imposta di bollo, tassa fissa durante il rogito.",
            },
        ];
    }

    function updateTable() {
        const tbody = document.getElementById("chartDataBody");
        if (!tbody) return;
        tbody.innerHTML = "";

        chartData.forEach((item, index) => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${item.name}</td>
                <td>${item.category}</td>
                <td>${item.value}</td>
                <td>${item.estimate ? "Yes" : "No"}</td>
            `;
            tbody.appendChild(row);
        });
    }
</script>

<main class="w-1/2 mx-auto p-6 border rounded-lg shadow flex flex-col gap-4">
    <h1 class="text-2xl font-bold text-center">🏡 House Cost Simulator</h1>

    <div class="flex flex-col gap-2">
        <label for="price">House Price (€)</label>
        <input
            id="price"
            type="number"
            bind:value={house_price}
            on:input={updateChart}
            class="w-full p-2 border rounded"
        />
        <input
            type="range"
            bind:value={house_price}
            min="0"
            max="800000"
            step="1000"
            on:input={updateChart}
            class="w-full"
        />
    </div>

    <label for="isFirstHouse">First house?</label>
    <input
        type="checkbox"
        id="firstHouse"
        bind:checked={is_fisrt_house}
        on:change={updateChart}
        class="w-5 h-5 border-2 rounded cursor-pointer"
    />

    <label for="isSoldByBuilder">Sold by builder within 5 years?</label>
    <input
        type="checkbox"
        id="isSoldByBuilder"
        bind:checked={is_sold_by_builder}
        on:change={updateChart}
        class="w-5 h-5 border-2 rounded cursor-pointer"
    />

    <label for="agencyFee">Agency Fee (%)</label>
    <input
        id="agencyFee"
        type="number"
        bind:value={agencyFee}
        on:input={updateChart}
        class="w-full p-2 border rounded"
    />

    <h2 class="text-xl font-semibold text-center">
        Total Cost: €{Number(total_cost).toLocaleString().split(".")[0]}
    </h2>

    <h2>Chart Data</h2>
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Category</th>
                <th>Value</th>
                <th>Estimate</th>
            </tr>
        </thead>
        <tbody>
            {#each chartData as item, index}
                <tr>
                    <td>{item.name}</td>
                    <td>{item.category}</td>
                    <td>
                        <input
                            type="number"
                            bind:value={chartData[index].value}
                        />
                    </td>
                    <td>{item.estimate ? "Yes" : "No"}</td>
                </tr>
            {/each}
        </tbody>
    </table>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        gap: 1rem; /* Adjust the gap size as needed */
    }
</style>
