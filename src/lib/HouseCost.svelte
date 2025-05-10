<script lang="ts">
    import { onMount, tick } from "svelte";
    import Chart from "chart.js/auto";
    import { fade, slide } from "svelte/transition";
    import Check from "../assets/Check.svelte";
    import NavBar from "../assets/NavBar.svelte";
    import DetailTile from "../assets/DetailTile.svelte";
    import TotalPriceTile from "../assets/TotalPriceTile.svelte";
    import SummaryDeltaPriceTile from "../assets/SummaryDeltaPriceTile.svelte";
    import DeltaPriceTile from "../assets/DeltaPriceTile.svelte";
    import TooltipFirstHouse from "../assets/TooltipFirstHouse.svelte";
    import TooltipVat from "../assets/TooltipVAT.svelte";
    import TooltipAgency from "../assets/TooltipAgency.svelte";
    import TooltipMortgage from "../assets/TooltipMortgage.svelte";
    import CustomButton from "../assets/CustomButton.svelte";

    let selectedTab = "summary";
    let prevSelectedTab = "";

    let house_price: number = 300000;
    let is_fisrt_house: boolean = false;
    let is_fisrt_house_label: string = "Agevolazione prima casa?";
    let is_sold_by_builder: boolean = false;
    let is_sold_by_builder_label: string = "Vendita soggetta a IVA?";
    let is_sold_by_agency: boolean = false;
    let is_sold_by_agency_label: string = "Acquisto tramite agenzia?";
    let is_using_mortgage: boolean = false;
    let is_using_mortgage_label: string = "Utilizzi un mutuo?";

    //mortgage variable section
    let taeg: number = 0;
    let mortgageDuration: number = 0;
    let mortgage_amount: number = 0;
    let interestsVal: number = 0;
    let mortgageInstallment: number = 0;
    let agencyFee: number = 0;
    let showCosts: boolean = false;
    let showTooltip: boolean[] = [false, false, false, false];

    let totalAmount: number = 0;
    let barChartInstance: Chart | null = null; // Store chart instance globally
    let interestsBarChart: Chart | null = null; // Store chart instance globally
    let dataByCategory: any = {};

    // Define category colors
    const categoryColors: Record<string, string> = {
        Tax: "rgba(133, 81, 182, 1)",
        Agency: "rgba(249, 166, 0, 1)",
        Notary: "rgba(98, 182, 170, 1)",
        Bank: "rgba(76, 51, 141, 1)",
    };
    const categoryBorderColors: Record<string, string> = {
        Tax: "rgba(76, 51, 141, 1)",
        Agency: "rgba(249, 166, 0, 1)",
        Notary: "rgba(98, 182, 170, 1)",
        Bank: "rgba(76, 51, 141, 1)",
    };

    // Step 1: Group data by category and sum values
    let groupedData: Record<string, number> = {};

    //updating the chart when switching nav
    $: if (
        selectedTab == "summary" &&
        prevSelectedTab != "summary" &&
        showCosts
    ) {
        prevSelectedTab = selectedTab;
        tick().then(() => {
            initializeDoughnutChart();
            updateData();
        });
    }

    $: if (
        selectedTab == "mortgage" &&
        prevSelectedTab != "mortgage" &&
        showCosts
    ) {
        prevSelectedTab = selectedTab;
        tick().then(() => {
            initializeInterestsBarChart();
            updateData();
        });
    }

    $: if (selectedTab == "base") {
        prevSelectedTab = selectedTab;
    }

    onMount(() => {});

    function buildApiString(): string {
        let apiStringUrl: string =
            "http://localhost:8080/get_house_prices?house_price=" +
            (house_price != null ? house_price : 0);

        if (is_sold_by_agency)
            apiStringUrl +=
                "&agency_fee=" + (agencyFee != null ? agencyFee / 100 : 0);
        if (is_fisrt_house)
            apiStringUrl +=
                "&is_first_house=" +
                (is_fisrt_house != null ? is_fisrt_house : false);
        if (is_sold_by_builder)
            apiStringUrl +=
                "&is_sold_by_builder=" +
                (is_sold_by_builder != null ? is_sold_by_builder : false);
        if (is_using_mortgage)
            apiStringUrl +=
                "&mortgage_amount=" +
                (mortgage_amount != null ? mortgage_amount : 0) +
                "&mortgage_duration=" +
                (mortgageDuration != null ? mortgageDuration : 0) +
                "&mortgage_TAEG=" +
                (taeg != null ? taeg / 100 : 0);

        return apiStringUrl;
    }

    async function updateData() {
        let chartData: any = [];

        const apiStringUrl: string = buildApiString();

        const apiCall = async () => {
            try {
                const response = await fetch(apiStringUrl);

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                chartData = await response.json();
                chartData = chartData["data"];
            } catch (error) {
                console.error("Fetch error:", error);
            }
        };

        apiCall().then(() => {
            dataByCategory = {};

            chartData.forEach((item: any) => {
                if (item.name == "Monthly Mortgage Installment") {
                    mortgageInstallment = item.value;
                }
                item.value = Math.round(item.value);

                //remove 0 values
                if (item.value > 0) {
                    if (item.category in dataByCategory) {
                        dataByCategory[item.category].push(item);
                    } else {
                        dataByCategory[item.category] = [item];
                    }
                }
            });
            groupedData = {};
            totalAmount = 0;

            chartData.forEach((item: any) => {
                if (
                    item.value > 0 &&
                    item.category != "House price" &&
                    item.category != "Installment"
                ) {
                    groupedData[item.category] =
                        (groupedData[item.category] || 0) +
                        Math.round(item.value);
                    totalAmount += item.value;
                }
            });
            updateDoughnutChart();
            updateInterestsBarChart();
        });
    }

    async function initializeDoughnutChart() {
        if (showCosts) {
            await tick(); // Wait for DOM update before selecting the canvas

            const ctx = document.getElementById(
                "barChart",
            ) as HTMLCanvasElement | null;
            if (!ctx) return; // Ensure ctx exists before using it

            const ctx2D = ctx.getContext("2d");
            if (!ctx2D) return; // Ensure context is available

            // Step 3: Destroy previous chart instance
            if (barChartInstance) {
                barChartInstance.destroy();
            }

            const categories = Object.keys(groupedData);
            const values = Object.values(groupedData);

            // Step 4: Create the stacked bar chart
            barChartInstance = new Chart(ctx2D, {
                type: "doughnut",
                data: {
                    labels: categories,
                    datasets: [
                        {
                            label: "Costo",
                            data: values,
                            backgroundColor: categories.map(
                                (cat) =>
                                    categoryColors[cat] ||
                                    "rgba(200, 200, 200, 0.6)",
                            ), // fallback color
                            borderWidth: 1,
                        },
                    ],
                },
                options: {
                    responsive: true,
                    plugins: {
                        legend: {
                            position: "bottom",
                            labels: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                },
                                padding: 20,
                            },
                        },
                        tooltip: {
                            mode: "index",
                            intersect: false,
                            callbacks: {
                                label: function (context) {
                                    const label = context.dataset.label || "";
                                    const value = context.parsed;
                                    const formatted = new Intl.NumberFormat(
                                        "it-IT",
                                        {
                                            style: "currency",
                                            currency: "EUR",
                                            minimumFractionDigits: 0,
                                            maximumFractionDigits: 0,
                                        },
                                    ).format(value);
                                    return `${label}: ${formatted}`;
                                },
                            },
                        },
                    },
                },
            });
        }
    }

    async function initializeInterestsBarChart() {
        //initializing interests bar chart

        if (showCosts) {
            const ictx = document.getElementById(
                "interestsBarChart",
            ) as HTMLCanvasElement | null;
            if (!ictx) return;

            const ictx2D = ictx.getContext("2d");
            if (!ictx2D) return;

            if (interestsBarChart) {
                interestsBarChart.destroy();
            }

            interestsBarChart = new Chart(ictx2D, {
                type: "bar",
                data: {
                    labels: [],
                    datasets: [
                        {
                            label: "Interests",
                            data: [],
                            backgroundColor: ["rgba(133, 81, 182, 1)"],
                            borderColor: ["rgba(133, 81, 182, 1)"],
                            borderWidth: 1,
                        },
                        {
                            label: "Capital",
                            data: [],
                            backgroundColor: ["rgba(249, 166, 0, 1)"],
                            borderColor: ["rgba(249, 166, 0, 1)"],
                            borderWidth: 1,
                        },
                    ],
                },
                options: {
                    scales: {
                        y: {
                            stacked: true,
                            beginAtZero: true,
                            min: 0,
                            max: mortgageInstallment * 12,
                            ticks: {
                                stepSize: mortgageInstallment * 12,
                                color: "white",
                                font: {
                                    weight: "bold",
                                },
                                callback: function (value) {
                                    // Only show the max value
                                    if (value === mortgageInstallment * 12) {
                                        return new Intl.NumberFormat("it-IT", {
                                            style: "currency",
                                            currency: "EUR",
                                            minimumFractionDigits: 0,
                                            maximumFractionDigits: 0,
                                        }).format(value);
                                    }
                                    return "";
                                },
                            },
                            grid: {
                                display: false,
                            },
                        },
                        x: {
                            stacked: true,
                            ticks: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                },
                            },
                            grid: {
                                display: false,
                            },
                        },
                    },
                    plugins: {
                        legend: {
                            labels: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                },
                                padding: 20,
                            },
                        },
                        tooltip: {
                            mode: "index",
                            intersect: false,
                            callbacks: {
                                label: function (context) {
                                    const label = context.dataset.label || "";
                                    const value = context.parsed.y;
                                    const formatted = new Intl.NumberFormat(
                                        "it-IT",
                                        {
                                            style: "currency",
                                            currency: "EUR",
                                            minimumFractionDigits: 0,
                                            maximumFractionDigits: 0,
                                        },
                                    ).format(value);
                                    return `${label}: ${formatted}`;
                                },
                            },
                        },
                    },
                },
            });
        }
    }

    function updateDoughnutChart() {
        if (!barChartInstance) return; // Ensure the chart exists

        // Step 2: Convert grouped data into datasets
        const categories = Object.keys(groupedData);
        const values = Object.values(groupedData);

        barChartInstance.data.labels = categories;

        barChartInstance.data.datasets[0].data = values; // Update dataset (values and colors)
        barChartInstance.data.datasets[0].backgroundColor = categories.map(
            (cat) => categoryColors[cat] || "rgba(200, 200, 200, 0.6)", // fallback
        );
        barChartInstance.data.datasets[0].borderColor = categories.map(
            (cat) => categoryBorderColors[cat] || "rgba(200, 200, 200, 0.6)", // fallback
        );
        barChartInstance.update(); // Smoothly update the chart
    }

    function updateInterestsBarChart() {
        if (!interestsBarChart) return; // Ensure the chart exists

        interestsVal = 0;
        const labels: string[] = [];
        const interestsData: number[] = [];
        const capitalData: number[] = [];
        let residualCapital = mortgage_amount;
        let interestsSum = 0;
        let capitalSum = 0;

        if (is_using_mortgage) {
            for (let i = 1; i <= mortgageDuration * 12; i++) {
                let interests: number = 0;
                let capital: number = 0;

                interests = (residualCapital * taeg) / 100 / 12;
                interestsSum += interests;
                interestsVal += interests;

                capital = mortgageInstallment - interests;
                capitalSum += capital;
                residualCapital -= capital;

                if (i % 12 == 0) {
                    labels.push(
                        (new Date().getFullYear() + (i / 12 - 1)).toString(),
                    );
                    interestsData.push(Math.round(interestsSum));
                    capitalData.push(Math.round(capitalSum));

                    interestsSum = 0;
                    capitalSum = 0;
                }
            }
        }

        interestsBarChart.data.labels = labels; // Update dataset (values and colors)
        interestsBarChart.data.datasets = [
            {
                label: "Interests",
                data: interestsData,
                backgroundColor: ["rgba(249, 166, 0, 0.9)"],
                borderColor: ["rgba(249, 166, 0, 1)"],
                borderWidth: 1,
            },
            {
                label: "Capital",
                data: capitalData,
                backgroundColor: ["rgba(133, 81, 182, 0.9)"],
                borderColor: ["rgba(133, 81, 182, 1)"],
                borderWidth: 1,
            },
        ];

        interestsBarChart.update(); // Smoothly update the chart
    }

    //added to upscale chart in case of window-resizing
    window.addEventListener("resize", () => {
        if (barChartInstance) {
            barChartInstance.resize(); // Force the chart to resize
        }
    });
</script>

<div
    class="h-screen bg-gradient-to-b from-purple-400 to-[#1e1f25] flex items-center justify-center p-6"
>
    <div
        class="relative flex w-[90%] max-w-8xl h-[80vh] transition-all duration-700 ease-in-out gap-6"
    >
        <!-- Left Tile: 40% -->
        <div
            class="basis-[40%] bg-[#1e1f25] text-white rounded-2xl shadow-lg p-8 space-y-6 overflow-auto transition-all duration-700 ease-in-out"
        >
            <h1 class="text-4xl font-bold leading-tight">
                Scopri quanto costa,<br />
                <span class="text-purple-400"> per davvero, </span><br />
                la casa dei tuoi sogni. 🏡
            </h1>
            <p class="text-gray-400 text-sm">
                Il prezzo finale della casa comprende tutte le spese accessorie
                per l'acquisto ed il totale delle tasse dovute. Per spese come
                il notaio è stata fatta una stima accurata su tutte le voci di
                costo.
            </p>

            <div class="flex flex-col gap-2 md:max-w-[200px]">
                <h1 class="text-base font-bold leading-tight">Costo casa:</h1>
                <div class="relative">
                    <input
                        type="number"
                        class="w-full border border-white rounded px-4 pr-10 py-2 text-left font-medium text-white"
                        min="0"
                        max="1000000"
                        step="5000"
                        bind:value={house_price}
                        on:change={showCosts ? updateData : () => {}}
                        on:mouseup={showCosts ? updateData : () => {}}
                    />
                    <!-- Separator Line -->
                    <div class="absolute inset-y-0 right-9 w-px bg-white"></div>
                    <!-- Euro Symbol -->
                    <div
                        class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                    >
                        €
                    </div>
                </div>
            </div>
            <div id="isFirstHouse">
                <Check
                    bind:label={is_fisrt_house_label}
                    bind:bind={is_fisrt_house}
                    bind:showTooltip
                    element={0}
                    on:change={showCosts ? updateData : () => {}}
                />
            </div>
            <div id="isSoldByBuilder">
                <Check
                    bind:label={is_sold_by_builder_label}
                    bind:bind={is_sold_by_builder}
                    bind:showTooltip
                    element={1}
                    on:change={showCosts ? updateData : () => {}}
                />
            </div>
            <div id="isSoldByAgency">
                <Check
                    bind:label={is_sold_by_agency_label}
                    bind:bind={is_sold_by_agency}
                    bind:showTooltip
                    element={2}
                    on:change={showCosts ? updateData : () => {}}
                />
                {#if is_sold_by_agency}
                    <div
                        class="flex flex-col gap-2 md:max-w-[200px]"
                        transition:slide={{ duration: 300 }}
                    >
                        <div class="ml-6">
                            <h1 class="text-base font-bold leading-tight mt-3">
                                Provvigione:
                            </h1>
                            <div class="relative">
                                <input
                                    type="number"
                                    class="w-full border border-white rounded px-4 pr-10 py-2 text-left font-medium text-white"
                                    min="0"
                                    max="7"
                                    step="0.1"
                                    bind:value={agencyFee}
                                    on:change={showCosts
                                        ? updateData
                                        : () => {}}
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                />
                                <!-- Separator Line -->
                                <div
                                    class="absolute inset-y-0 right-9 w-px bg-white"
                                ></div>
                                <!-- % Symbol -->
                                <div
                                    class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                >
                                    %
                                </div>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>
            <div id="isUsingMortgage">
                <Check
                    bind:label={is_using_mortgage_label}
                    bind:bind={is_using_mortgage}
                    bind:showTooltip
                    element={3}
                    on:change={showCosts ? updateData : () => {}}
                />
                {#if is_using_mortgage}
                    <div
                        class="flex flex-col gap-2 md:max-w-[200px]"
                        transition:slide={{ duration: 300 }}
                    >
                        <div class="ml-6">
                            <h1 class="text-base font-bold leading-tight mt-3">
                                Importo mutuo:
                            </h1>
                            <div class="relative">
                                <input
                                    type="number"
                                    class="w-full border border-white rounded px-4 pr-10 py-2 text-left font-medium text-white"
                                    min="0"
                                    max={house_price}
                                    step="5000"
                                    bind:value={mortgage_amount}
                                    on:change={showCosts
                                        ? updateData
                                        : () => {}}
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                />
                                <!-- Separator Line -->
                                <div
                                    class="absolute inset-y-0 right-9 w-px bg-white"
                                ></div>
                                <!-- % Symbol -->
                                <div
                                    class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                >
                                    €
                                </div>
                            </div>
                        </div>

                        <div class="ml-6">
                            <h1 class="text-base font-bold leading-tight mt-3">
                                Durata (Anni):
                            </h1>
                            <div class="relative">
                                <input
                                    type="number"
                                    class="w-full border border-white rounded px-4 pr-10 py-2 text-left font-medium text-white"
                                    min="0"
                                    bind:value={mortgageDuration}
                                    on:change={showCosts
                                        ? updateData
                                        : () => {}}
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                />
                                <!-- Separator Line -->
                                <div
                                    class="absolute inset-y-0 right-9 w-px bg-white"
                                ></div>
                                <!-- % Symbol -->
                                <div
                                    class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                >
                                    €
                                </div>
                            </div>
                        </div>

                        <div class="ml-6">
                            <h1 class="text-base font-bold leading-tight mt-3">
                                TAEG:
                            </h1>
                            <div class="relative">
                                <input
                                    type="number"
                                    class="w-full border border-white rounded px-4 pr-10 py-2 text-left font-medium text-white"
                                    min="0"
                                    step="0.1"
                                    bind:value={taeg}
                                    on:change={showCosts
                                        ? updateData
                                        : () => {}}
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                />
                                <!-- Separator Line -->
                                <div
                                    class="absolute inset-y-0 right-9 w-px bg-white"
                                ></div>
                                <!-- % Symbol -->
                                <div
                                    class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                >
                                    %
                                </div>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>
            {#if !showCosts}
                <div
                    transition:slide={{ duration: 300 }}
                    class="overflow-hidden"
                >
                    <div
                        class="flex items-center space-x-2 transition-all duration-300"
                    >
                        <CustomButton
                            title={"Quanto mi costi?"}
                            onClick={() => {
                                for (let i = 0; i < showTooltip.length; i++) {
                                    showTooltip[i] = false;
                                }
                                showCosts = true;
                            }}
                        ></CustomButton>
                    </div>
                </div>
            {/if}
        </div>
        <!-- Right Tile + Tooltip Wrapper -->
        <div class="relative basis-[60%]">
            <!-- Right Tile: 60% -->
            <div
                class="bg-[#1e1f25] rounded-2xl shadow-lg p-8 overflow-auto w-full h-full transition-all duration-600 ease-out-[cubic-bezier(0.22, 1, 0.36, 1)]"
                class:opacity-100={showCosts}
                class:opacity-0={!showCosts}
                class:-translate-y-500={!showCosts}
                class:translate-y-0={showCosts}
                class:pointer-events-none={!showCosts}
            >
                <div class="max-w-6xl mx-auto">
                    <!-- Costs & Chart Section -->
                    {#if showCosts}
                        <div transition:fade={{ duration: 500 }}>
                            <NavBar bind:selectedTab {is_using_mortgage} />
                            {#if selectedTab == "summary"}
                                <!-- Main Chart + Cost Breakdown Side-by-Side -->
                                <div transition:fade={{ duration: 500 }}>
                                    <div
                                        class="flex flex-wrap gap-8 justify-center items-center w-full h-full"
                                        transition:slide={{ duration: 500 }}
                                    >
                                        <!-- Chart -->
                                        <div
                                            class="flex justify-center w-full max-w-70 h-[150px] sm:h-[200px] md:h-[250px] lg:h-[300px]"
                                        >
                                            <canvas id="barChart"></canvas>
                                        </div>
                                        <!-- Total (Main) Price -->
                                        <TotalPriceTile
                                            number={totalAmount + house_price}
                                            dif={totalAmount}
                                            name={"💰 Spesa totale"}
                                        />
                                    </div>

                                    <div
                                        class="flex flex-wrap gap-4 mt-8 w-full"
                                        transition:slide={{ duration: 500 }}
                                    >
                                        <SummaryDeltaPriceTile
                                            name={"🇮🇹 Tasse"}
                                            number={groupedData["Tax"]}
                                            showVal={house_price != 0}
                                            delta={groupedData["Tax"] /
                                                house_price}
                                        />
                                        <!-- Subtiles Container with Vertical Line -->

                                        <SummaryDeltaPriceTile
                                            name={"📜 Notaio"}
                                            number={groupedData["Notary"]}
                                            showVal={house_price != 0}
                                            delta={groupedData["Notary"] /
                                                house_price}
                                        />
                                        {#if is_sold_by_agency}
                                            <!-- Agency -->
                                            <SummaryDeltaPriceTile
                                                name={"🏘️ Agenzia"}
                                                number={groupedData["Agency"]}
                                                showVal={house_price > 0 &&
                                                    groupedData["Agency"] !=
                                                        null}
                                                delta={groupedData["Agency"] /
                                                    house_price}
                                            />
                                        {/if}
                                        {#if is_using_mortgage}
                                            <!-- Bank Costs -->
                                            <SummaryDeltaPriceTile
                                                name={"🏦 Banca"}
                                                number={groupedData["Bank"]}
                                                showVal={groupedData["Bank"] !=
                                                    null}
                                                delta={groupedData["Bank"] /
                                                    house_price}
                                            />
                                        {/if}
                                    </div>
                                </div>
                            {/if}
                            {#if selectedTab == "base"}
                                <div transition:fade={{ duration: 500 }}>
                                    <div
                                        class="flex flex-wrap gap-1"
                                        transition:slide={{ duration: 500 }}
                                    >
                                        <DetailTile
                                            data={dataByCategory}
                                            element={"Tax"}
                                            title={"🇮🇹 Taxes"}
                                        />
                                        <DetailTile
                                            data={dataByCategory}
                                            element={"Notary"}
                                            title={"📜 Notary"}
                                        />

                                        {#if "Agency" in dataByCategory}
                                            <DetailTile
                                                data={dataByCategory}
                                                element={"Agency"}
                                                title={"🏘️ Agency"}
                                            />
                                        {/if}
                                        {#if "Bank" in dataByCategory}
                                            <DetailTile
                                                data={dataByCategory}
                                                element={"Bank"}
                                                title={"🏦 Bank"}
                                            />
                                        {/if}
                                    </div>
                                </div>
                            {/if}
                            {#if selectedTab == "mortgage"}
                                <!-- Entire mortgage section -->
                                <section
                                    class="flex flex-col w-full gap-6"
                                    transition:slide={{ duration: 500 }}
                                >
                                    <div transition:fade={{ duration: 500 }}>
                                        <div
                                            class="flex flex-wrap gap-6"
                                            transition:slide={{
                                                duration: 500,
                                            }}
                                        >
                                            <div
                                                class="w-full md:w-[48%] mb-6"
                                                transition:slide={{
                                                    duration: 500,
                                                }}
                                            >
                                                <SummaryDeltaPriceTile
                                                    number={mortgageInstallment}
                                                    name={"Rata mutuo"}
                                                    showVal={mortgageInstallment !=
                                                        null}
                                                    delta={null}
                                                />
                                            </div>
                                            <div
                                                class="w-full md:w-[48%] mb-6"
                                                transition:slide={{
                                                    duration: 500,
                                                }}
                                            >
                                                <DeltaPriceTile
                                                    number={is_using_mortgage
                                                        ? interestsVal
                                                        : 0}
                                                    name={"Interests in"}
                                                    duration={mortgageDuration}
                                                    delta={interestsVal /
                                                        mortgage_amount}
                                                    showVal={mortgage_amount !=
                                                        null &&
                                                        mortgage_amount != 0 &&
                                                        is_using_mortgage}
                                                />
                                            </div>
                                        </div>
                                        <!-- Chart -->
                                        <div
                                            class="flex items-center justify-center"
                                            transition:slide={{
                                                duration: 500,
                                            }}
                                        >
                                            <canvas
                                                id="interestsBarChart"
                                                class="w-full h-full"
                                            ></canvas>
                                        </div>
                                    </div>
                                </section>
                            {/if}
                        </div>
                    {/if}
                </div>
            </div>
            <!-- Tooltip (Overlay) -->
            <TooltipFirstHouse showTooltip={showTooltip[0]} />
            <TooltipVat showTooltip={showTooltip[1]} />
            <TooltipAgency showTooltip={showTooltip[2]} />
            <TooltipMortgage showTooltip={showTooltip[3]} />
        </div>
    </div>
</div>

<style>
</style>
