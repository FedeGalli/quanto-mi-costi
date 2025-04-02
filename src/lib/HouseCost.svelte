<script lang="ts">
    import { onMount, tick } from "svelte";
    import Chart from "chart.js/auto";
    import { fade, scale } from "svelte/transition";
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    import Check from "../assets/check.svelte";

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
    let is_fisrt_house: boolean = false;
    let is_fisrt_house_label: string = "First House?";
    let is_sold_by_builder: boolean = false;
    let is_sold_by_builder_label: string = "Sold by builder within 5 years?";
    let is_sold_by_agency: boolean = false;
    let is_sold_by_agency_label: string = "Agency?";

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

    let agencyFee: number = 0;
    let agencyAmount: number = agencyFee * house_price;
    let notaryAmount: number = house_price * 0.00085 * 1.22;

    let imposta_di_bollo = 230;
    let tassa_ipotecaria = 90;
    let tassa_archivio = 60;
    let visure_ipotecarie = 180;
    let showCosts = false;

    //if the number is < 0 i assume is the integer one
    $: house_price = house_price >= 0 ? house_price : house_price * -1;
    $: agencyFee = agencyFee >= 0 ? agencyFee : agencyFee * -1;

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

    $: agencyAmount = is_sold_by_agency ? agencyFee * house_price * 1.22 : 0;

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
    let doughnutChartInstance: Chart | null = null; // Store chart instance globally

    onMount(() => {});

    function updateData() {
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

        //charts
        updateDoughnut(chartData);
    }

    async function initializeDoughnut() {
        if (showCosts) {
            await tick(); // Wait for DOM update before selecting the canvas

            const ctx = document.getElementById("myChart") as HTMLCanvasElement;

            if (!ctx) return; // Ensure ctx exists before using it
            const ctx2D = ctx.getContext("2d"); // Get the 2D drawing context

            if (!ctx2D) return; // Ensure context is available

            const filteredData = chartData.filter(
                (item) => item.value > 0 && item.name != "Property Price",
            ); //filtering out 0 values

            doughnutChartInstance = new Chart(ctx2D, {
                type: "doughnut",
                data: {
                    labels: filteredData.map((item) => item.name), // Names as labels
                    datasets: [
                        {
                            label: "Cost Values",
                            data: filteredData.map((item) => item.value), // Values as data
                            backgroundColor: filteredData.map(
                                (item) =>
                                    item.category == "Tax"
                                        ? "rgba(255, 99, 132, 0.6)"
                                        : item.category == "Agency"
                                          ? "rgba(54, 162, 235, 0.6)"
                                          : "rgba(201, 200, 115, 0.6)", //notary category
                            ),
                            borderColor: filteredData.map(
                                (item) =>
                                    item.category == "Tax"
                                        ? "rgba(255, 99, 132)"
                                        : item.category == "Agency"
                                          ? "rgba(54, 162, 235)"
                                          : "rgba(201, 200, 115)", //notary category
                            ),
                            borderWidth: 1,
                            //borderJoinStyle: "round",
                            //borderRadius: 100,
                        },
                    ],
                },

                options: {
                    responsive: true,
                    animation: {
                        animateRotate: true, // Enable rotation animation
                        animateScale: true, // Optional: Animates scaling of the pie chart
                        duration: 1000, // Duration of animation (1 seconds)
                        easing: "easeInOutCubic", // Smooth easing effect
                    },
                },
            });
        }
    }

    function updateDoughnut(newData: CostItem[]) {
        if (!doughnutChartInstance) return; // Ensure the chart exists

        const filteredData = chartData.filter(
            (item) => item.value > 0 && item.name != "Property Price",
        ); //filtering out 0 values

        doughnutChartInstance.data.labels = filteredData.map(
            (item) => item.name,
        ); // Update labels
        doughnutChartInstance.data.datasets[0].data = filteredData.map(
            (item) => item.value,
        ); // Update values
        doughnutChartInstance.data.datasets[0].backgroundColor =
            filteredData.map(
                (item) =>
                    item.category == "Tax"
                        ? "rgba(255, 99, 132,0.6)"
                        : item.category == "Agency"
                          ? "rgba(54, 162, 235, 0.6)"
                          : "rgba(201, 200, 115, 0.6)", //notary category
            ); // Update backgroud colors
        doughnutChartInstance.data.datasets[0].borderColor = filteredData.map(
            (item) =>
                item.category == "Tax"
                    ? "rgba(255, 99, 132)"
                    : item.category == "Agency"
                      ? "rgba(54, 162, 235)"
                      : "rgba(201, 200, 115)", //notary category
        ); // Update border colors

        doughnutChartInstance.update(); // Smoothly update the chart
    }

    function calculateCosts() {
        showCosts = true;
        initializeDoughnut();
        updateData();
    }

    //added to upscale chart in case of window-resizing
    window.addEventListener("resize", () => {
        if (doughnutChartInstance) {
            doughnutChartInstance.resize(); // Force the chart to resize
        }
    });
</script>

<main class="w-full p-6 flex flex-col gap-8 items-center">
    <!-- Inputs Section -->
    <section
        class="w-full max-w-3xl p-6 border rounded-lg shadow flex flex-col gap-6"
        transition:scale={{ duration: 500 }}
    >
        <h1 class="text-2xl font-bold text-center">
            🏡 Real House Cost Simulator
        </h1>

        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2 md:max-w-[360px]">
                <label for="price"
                    >House Price: <NumberFlow
                        value={house_price}
                        format={{
                            style: "currency",
                            currency: "EUR",
                            maximumFractionDigits: 0,
                        }}
                        locales={"it-IT"}
                    /></label
                >
                <input
                    type="range"
                    bind:value={house_price}
                    min="0"
                    max="800000"
                    step="1000"
                    on:input={updateData}
                    class="w-full"
                />
            </div>
            <div id="isFirstHouse">
                <Check
                    bind:label={is_fisrt_house_label}
                    bind:bind={is_fisrt_house}
                    on:change={updateData}
                />
            </div>
            <div id="isSoldByBuilder">
                <Check
                    bind:label={is_sold_by_builder_label}
                    bind:bind={is_sold_by_builder}
                    on:change={updateData}
                />
            </div>
            <div id="isSoldByAgency">
                <Check
                    bind:label={is_sold_by_agency_label}
                    bind:bind={is_sold_by_agency}
                    on:change={updateData}
                />
                {#if is_sold_by_agency}
                    <div class="flex flex-col gap-2 md:max-w-[360px]">
                        <label for="Fee"
                            ><NumberFlow
                                value={agencyFee}
                                format={{
                                    style: "percent",
                                    minimumFractionDigits: 2,
                                }}
                                locales={"it-IT"}
                            /></label
                        >
                        <input
                            type="range"
                            bind:value={agencyFee}
                            min="0"
                            max="0.1"
                            step="0.001"
                            on:input={updateData}
                            class="w-full"
                        />
                    </div>
                {/if}
            </div>
        </div>

        {#if !showCosts}
            <div class="flex items-center space-x-2">
                <button on:click={calculateCosts} class="btn">
                    {showCosts ? "Hide Costs" : "Quanto mi costi?"}
                </button>
            </div>
        {/if}
    </section>

    <!-- Costs & Chart Section -->
    {#if showCosts}
        <section
            class="flex flex-col w-full max-w-3xl gap-6"
            transition:fade={{ duration: 2000 }}
        >
            <!-- Cost Summary -->
            <div class="flex flex-col gap-3">
                <div class="text-xl flex-row gap-3 font-semibold text-center">
                    <NumberFlowGroup
                        style="--number-flow-char-height: 0.85em"
                        class="flex items-center gap-4 font-semibold"
                    >
                        <NumberFlow
                            value={house_price}
                            format={{
                                style: "currency",
                                currency: "EUR",
                                maximumFractionDigits: 0,
                            }}
                            locales={"it-IT"}
                        />
                        <NumberFlow
                            value={total_cost}
                            format={{
                                style: "currency",
                                currency: "EUR",
                                maximumFractionDigits: 0,
                            }}
                            class="text-3xl"
                        />
                        <NumberFlow
                            value={total_cost / house_price - 1}
                            format={{
                                style: "percent",
                                maximumFractionDigits: 2,
                                signDisplay: "always",
                            }}
                            class="text-lg transition-colors duration-300 text-red-500"
                        />
                    </NumberFlowGroup>
                </div>
                {#each chartData as item, index}
                    <h1>
                        <b>{item.name}</b>
                    </h1>
                    <NumberFlow
                        value={chartData[index].value}
                        format={{
                            style: "currency",
                            currency: "EUR",
                            maximumFractionDigits: 0,
                        }}
                        locales={"it-IT"}
                    />
                {/each}
            </div>

            <!-- Chart -->
            <div class="border p-6 rounded-lg shadow w-full">
                <canvas id="myChart"></canvas>
            </div>
        </section>
    {/if}
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        gap: 2rem;
        align-items: center;
    }

    .btn {
        background-color: rgba(54, 158, 82, 0.8);
        color: white;
        padding: 10px 20px;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        font-size: 16px;
        transition: background-color 0.3s;
    }

    .btn:hover {
        background-color: rgba(54, 158, 82, 1);
    }
</style>
