<script lang="ts">
    import { onMount, tick } from "svelte";
    import Chart from "chart.js/auto";
    import { fade, scale, slide } from "svelte/transition";
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    import Check from "../assets/Check.svelte";
    import NavBar from "../assets/NavBar.svelte";

    // Define types for the house cost breakdown
    type CostItem = {
        name: string;
        category: string;
        value: number;
        estimate: boolean;
        info: string;
    };

    let selectedTab = "summary";

    let house_price: number = 300000;
    let is_fisrt_house: boolean = false;
    let is_fisrt_house_label: string = "First house?";
    let is_sold_by_builder: boolean = false;
    let is_sold_by_builder_label: string = "Sold by builder within 5 years?";
    let is_sold_by_agency: boolean = false;
    let is_sold_by_agency_label: string = "Agency?";
    let is_using_mortgage: boolean = false;
    let is_using_mortgage_label: string = "Using mortgage?";

    //mortgage variable section
    let taeg: number = 0;
    let mortgageDuration: number = 0;
    let mortgage_amount: number = 0;
    let interestsVal: number = 0;
    let mortgageInstallment: number = 0;
    let agencyFee: number = 0;
    let showCosts = false;

    let totalAmount: number = 0;
    let barChartInstance: Chart | null = null; // Store chart instance globally

    $: if (selectedTab == "summary") initializeBarChart();

    // Define category colors
    const categoryColors: Record<string, string> = {
        Tax: "rgba(255, 99, 132, 0.6)",
        Agency: "rgba(54, 162, 235, 0.6)",
        Notary: "rgba(201, 200, 115, 0.6)",
        Interests: "rgba(152, 3, 252, 0.6)",
        Bank: "rgba(66, 245, 126, 0.6)",
    };

    // Step 1: Group data by category and sum values
    let groupedData: Record<string, number> = {};

    onMount(() => {});

    function buildApiString(): string {
        let apiStringUrl: string =
            "http://localhost:8080/get_house_prices?house_price=" + house_price;

        if (is_sold_by_agency) apiStringUrl += "&agency_fee=" + agencyFee;
        if (is_fisrt_house) apiStringUrl += "&is_first_house=" + is_fisrt_house;
        if (is_sold_by_builder)
            apiStringUrl += "&is_sold_by_builder=" + is_sold_by_builder;
        if (is_using_mortgage)
            apiStringUrl +=
                "&mortgage_amount=" +
                mortgage_amount +
                "&mortgage_duration=" +
                mortgageDuration +
                "&mortgage_TAEG=" +
                taeg;

        return apiStringUrl;
    }

    async function updateData() {
        let chartData = [];

        const apiStringUrl: string = buildApiString();

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

        chartData.forEach((item: any) => {
            if (item.name == "Mortgage Interests") interestsVal = item.value;
        });

        mortgageInstallment =
            is_using_mortgage && mortgageDuration != 0
                ? interestsVal / mortgageDuration / 12 +
                  mortgage_amount / mortgageDuration / 12
                : 0;

        groupedData = {};
        totalAmount = 0;

        chartData.forEach((item: any) => {
            if (item.value > 0 && item.category != "House price") {
                groupedData[item.category] =
                    (groupedData[item.category] || 0) + Math.round(item.value);
            }
            totalAmount += item.value;
        });
        //charts
        updateBarChart(chartData);
    }

    async function initializeBarChart() {
        if (showCosts) {
            await tick(); // Wait for DOM update before selecting the canvas

            const ctx = document.getElementById(
                "myChart",
            ) as HTMLCanvasElement | null;
            if (!ctx) return; // Ensure ctx exists before using it

            const ctx2D = ctx.getContext("2d");
            if (!ctx2D) return; // Ensure context is available

            // Step 2: Convert grouped data into datasets
            const categories = Object.keys(groupedData);
            const values = Object.values(groupedData);

            const datasets = categories.map((category, index) => ({
                label: category,
                data: [values[index]], // Single bar with stacked values
                backgroundColor:
                    categoryColors[category] || "rgba(100, 100, 100, 0.6)",
                borderColor:
                    categoryColors[category]?.replace("0.6", "1") ||
                    "rgba(100, 100, 100, 1)",
                borderWidth: 1,
            }));

            // Step 3: Destroy previous chart instance
            if (barChartInstance) {
                barChartInstance.destroy();
            }

            const max = Math.round(totalAmount - house_price);

            // Step 4: Create the stacked bar chart
            barChartInstance = new Chart(ctx2D, {
                type: "bar",
                data: {
                    labels: ["Total Costs"],
                    datasets: datasets,
                },
                options: {
                    animation: {
                        duration: 1500,
                        easing: "easeOutBounce",
                    },
                    indexAxis: "x",
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        x: {
                            stacked: true,
                            grid: { display: false },
                        },
                        y: {
                            stacked: true,
                            grid: { display: false },
                            afterBuildTicks: function (axis) {
                                const ticks = [0];
                                axis.ticks = ticks.map((v) => ({ value: v }));
                            },
                        },
                    },
                },
            });
        }
    }

    function updateBarChart(newData: CostItem[]) {
        if (!barChartInstance) return; // Ensure the chart exists

        // Step 2: Convert grouped data into datasets
        const categories = Object.keys(groupedData);
        const values = Object.values(groupedData);

        const datasets = categories.map((category, index) => ({
            label: category,
            data: [values[index]], // Single bar with stacked values
            backgroundColor:
                categoryColors[category] || "rgba(100, 100, 100, 0.6)",
            borderColor:
                categoryColors[category]?.replace("0.6", "1") ||
                "rgba(100, 100, 100, 1)",
            borderWidth: 1,
        }));
        barChartInstance.data.datasets = datasets; // Update dataset (values and colors)

        let tickList: number[] = [];
        tickList[0] = values[0];

        values.forEach((item, i) => {
            if (i != 0) tickList[i] = item + tickList[i - 1];
        });

        //adjusting the max bar size to the maximum
        const max = Math.round(totalAmount - house_price);

        barChartInstance!.options!.scales!.y!.max = max;
        barChartInstance!.options!.scales!.y!.afterBuildTicks = function (
            axis,
        ) {
            const ticks = tickList;
            axis.ticks = ticks.map((v) => ({ value: v }));
        };
        barChartInstance.update(); // Smoothly update the chart
    }

    function calculateCosts() {
        showCosts = true;
        initializeBarChart();
        updateData();
    }

    //added to upscale chart in case of window-resizing
    window.addEventListener("resize", () => {
        if (barChartInstance) {
            barChartInstance.resize(); // Force the chart to resize
        }
    });
</script>

<div class="h-screen overflow-hidden">
    <div class="flex h-full">
        <div class="w-86 bg-gray-700 text-white p-4 flex flex-col">
            <h1 class="text-xl font-bold mb-6">🏡 Real House Cost Simulator</h1>
            <!-- Inputs Section -->
            <section
                class="w-full max-w-3xl p-6 flex flex-col gap-6"
                transition:fade={{ duration: 1000 }}
            >
                <div class="flex flex-col gap-4 mb-4">
                    <div class="flex flex-col gap-2 md:max-w-[360px]">
                        <label for="price"
                            ><b>House Price:</b>
                            <NumberFlow
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
                            type="number"
                            bind:value={house_price}
                            min="0"
                            max="1000000"
                            step="1000"
                            on:change={showCosts ? updateData : () => {}}
                            on:mouseup={showCosts ? updateData : () => {}}
                            class="border p-2 rounded w-full"
                        />
                        <input
                            type="range"
                            bind:value={house_price}
                            min="0"
                            max="1000000"
                            step="1000"
                            on:mouseup={showCosts ? updateData : () => {}}
                            class="w-full"
                        />
                    </div>
                    <div id="isFirstHouse">
                        <Check
                            bind:label={is_fisrt_house_label}
                            bind:bind={is_fisrt_house}
                            on:change={showCosts ? updateData : () => {}}
                        />
                    </div>
                    <div id="isSoldByBuilder">
                        <Check
                            bind:label={is_sold_by_builder_label}
                            bind:bind={is_sold_by_builder}
                            on:change={showCosts ? updateData : () => {}}
                        />
                    </div>
                    <div id="isSoldByAgency">
                        <Check
                            bind:label={is_sold_by_agency_label}
                            bind:bind={is_sold_by_agency}
                            on:change={showCosts ? updateData : () => {}}
                        />
                        {#if is_sold_by_agency}
                            <div
                                class="flex flex-col gap-2 md:max-w-[360px]"
                                transition:slide={{ duration: 300 }}
                            >
                                <label for="Fee"
                                    >Fee: <NumberFlow
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
                                    max="0.07"
                                    step="0.0005"
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                    class="w-full"
                                />
                            </div>
                        {/if}
                    </div>
                    <div id="isUsingMortgage">
                        <Check
                            bind:label={is_using_mortgage_label}
                            bind:bind={is_using_mortgage}
                            on:change={showCosts ? updateData : () => {}}
                        />
                        {#if is_using_mortgage}
                            <div
                                class="flex flex-col gap-2 md:max-w-[360px]"
                                transition:slide={{ duration: 300 }}
                            >
                                <label for="mortgage_amount"
                                    >Mortgage Amount: <NumberFlow
                                        value={mortgage_amount}
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
                                    bind:value={mortgage_amount}
                                    min="0"
                                    max={house_price}
                                    step="5000"
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                    class="w-full"
                                />
                                <label for="duration"
                                    >Duration(Years): <NumberFlow
                                        value={mortgageDuration}
                                        format={{
                                            minimumFractionDigits: 0,
                                        }}
                                        locales={"it-IT"}
                                    /></label
                                >
                                <input
                                    type="range"
                                    bind:value={mortgageDuration}
                                    min="5"
                                    max="40"
                                    step="5"
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                    class="w-full"
                                />
                                <label for="TAEG"
                                    >TAEG: <NumberFlow
                                        value={taeg}
                                        format={{
                                            style: "percent",
                                            minimumFractionDigits: 2,
                                        }}
                                        locales={"it-IT"}
                                    /></label
                                >
                                <input
                                    type="range"
                                    bind:value={taeg}
                                    min="0"
                                    max="0.08"
                                    step="0.0005"
                                    on:mouseup={showCosts
                                        ? updateData
                                        : () => {}}
                                    class="w-full"
                                />
                            </div>
                        {/if}
                    </div>
                </div>
                {#if !showCosts}
                    <div
                        transition:slide={{ duration: 300 }}
                        class="overflow-hidden"
                    >
                        <div
                            class="flex items-center space-x-2 transition-all duration-300"
                        >
                            <button on:click={calculateCosts} class="btn">
                                {showCosts ? "Hide Costs" : "Quanto mi costi?"}
                            </button>
                        </div>
                    </div>
                {/if}
            </section>
            <div class="mt-auto pt-4 border-t border-gray-500">
                <p class="text-sm text-gray-400">Footer content</p>
            </div>
        </div>

        <!-- Main content -->
        <div class="flex-1 flex flex-col overflow-hidden">
            <main class="flex-1 overflow-y-auto p-4">
                <div class="max-w-4xl mx-auto">
                    <!-- Costs & Chart Section -->
                    {#if showCosts}
                        <div transition:fade={{ duration: 2000 }}>
                            <NavBar bind:selectedTab />
                            {#if selectedTab == "summary"}
                                <!-- Entire summary section -->
                                <section
                                    class="flex flex-col w-full gap-6"
                                    transition:slide={{ duration: 500 }}
                                >
                                    <!-- Final Total Price on Top -->
                                    <div class="flex justify-center">
                                        <div
                                            class="border p-6 rounded-lg shadow w-full max-w-md text-center"
                                        >
                                            <h1 class="text-xl font-bold mb-2">
                                                💸 Final Total Price
                                            </h1>
                                            <NumberFlowGroup
                                                style="--number-flow-char-height: 0.85em"
                                                class="flex items-center gap-4 font-semibold"
                                            >
                                                <NumberFlow
                                                    value={totalAmount}
                                                    format={{
                                                        style: "currency",
                                                        currency: "EUR",
                                                        maximumFractionDigits: 0,
                                                    }}
                                                    locales={"it-IT"}
                                                    class="text-3xl"
                                                />
                                                <NumberFlow
                                                    value={totalAmount -
                                                        house_price}
                                                    format={{
                                                        style: "currency",
                                                        currency: "EUR",
                                                        signDisplay: "always",
                                                        maximumFractionDigits: 0,
                                                    }}
                                                    class="text-lg transition-colors duration-300 text-red-500"
                                                />
                                            </NumberFlowGroup>
                                        </div>
                                    </div>

                                    <!-- Main Chart + Cost Breakdown Side-by-Side -->
                                    <div
                                        class="flex flex-row gap-6 justify-center"
                                    >
                                        <!-- Chart -->
                                        <div
                                            class="border p-6 rounded-lg shadow w-full"
                                        >
                                            <canvas id="myChart" height="400"
                                            ></canvas>
                                        </div>

                                        <!-- Cost Breakdown Column -->
                                        <div
                                            class="flex flex-col gap-4 max-w-sm w-full"
                                        >
                                            <!-- Taxes -->
                                            <div
                                                class="border p-6 rounded-lg shadow text-center"
                                            >
                                                <h1
                                                    class="text-lg font-semibold mb-2"
                                                >
                                                    🇮🇹 Taxes
                                                </h1>
                                                <NumberFlowGroup
                                                    style="--number-flow-char-height: 0.85em"
                                                    class="flex items-center gap-4 font-semibold"
                                                >
                                                    <NumberFlow
                                                        value={groupedData[
                                                            "Tax"
                                                        ]}
                                                        format={{
                                                            style: "currency",
                                                            currency: "EUR",
                                                            maximumFractionDigits: 0,
                                                        }}
                                                        locales={"it-IT"}
                                                        class="text-3xl"
                                                    />
                                                    <NumberFlow
                                                        value={house_price == 0
                                                            ? 0
                                                            : groupedData[
                                                                  "Tax"
                                                              ] / house_price}
                                                        format={{
                                                            style: "percent",
                                                            maximumFractionDigits: 2,
                                                            minimumFractionDigits: 2,
                                                            signDisplay:
                                                                "always",
                                                        }}
                                                        class="text-lg transition-colors duration-300 text-red-500"
                                                    />
                                                </NumberFlowGroup>
                                            </div>

                                            <!-- Notary -->
                                            <div
                                                class="border p-6 rounded-lg shadow text-center"
                                            >
                                                <h1
                                                    class="text-lg font-semibold mb-2"
                                                >
                                                    💼 Notary
                                                </h1>
                                                <NumberFlowGroup
                                                    style="--number-flow-char-height: 0.85em"
                                                    class="flex items-center gap-4 font-semibold"
                                                >
                                                    <NumberFlow
                                                        value={groupedData[
                                                            "Notary"
                                                        ]}
                                                        format={{
                                                            style: "currency",
                                                            currency: "EUR",
                                                            maximumFractionDigits: 0,
                                                        }}
                                                        locales={"it-IT"}
                                                        class="text-3xl"
                                                    />
                                                    <NumberFlow
                                                        value={house_price == 0
                                                            ? 0
                                                            : groupedData[
                                                                  "Notary"
                                                              ] / house_price}
                                                        format={{
                                                            style: "percent",
                                                            maximumFractionDigits: 2,
                                                            minimumFractionDigits: 2,
                                                            signDisplay:
                                                                "always",
                                                        }}
                                                        class="text-lg transition-colors duration-300 text-red-500"
                                                    />
                                                </NumberFlowGroup>
                                            </div>

                                            {#if is_sold_by_agency}
                                                <!-- Agency -->
                                                <div
                                                    class="border p-6 rounded-lg shadow text-center"
                                                >
                                                    <h1
                                                        class="text-lg font-semibold mb-2"
                                                    >
                                                        🏘️ Agency
                                                    </h1>
                                                    <NumberFlowGroup
                                                        style="--number-flow-char-height: 0.85em"
                                                        class="flex items-center gap-4 font-semibold"
                                                    >
                                                        <NumberFlow
                                                            value={groupedData[
                                                                "Agency"
                                                            ]
                                                                ? groupedData[
                                                                      "Agency"
                                                                  ]
                                                                : 0}
                                                            format={{
                                                                style: "currency",
                                                                currency: "EUR",
                                                                maximumFractionDigits: 0,
                                                            }}
                                                            locales={"it-IT"}
                                                            class="text-3xl"
                                                        />
                                                        <NumberFlow
                                                            value={house_price <=
                                                                0 ||
                                                            !groupedData[
                                                                "Agency"
                                                            ]
                                                                ? 0
                                                                : groupedData[
                                                                      "Agency"
                                                                  ] /
                                                                  house_price}
                                                            format={{
                                                                style: "percent",
                                                                maximumFractionDigits: 2,
                                                                minimumFractionDigits: 2,
                                                                signDisplay:
                                                                    "always",
                                                            }}
                                                            class="text-lg transition-colors duration-300 text-red-500"
                                                        />
                                                    </NumberFlowGroup>
                                                </div>
                                            {/if}

                                            {#if is_using_mortgage}
                                                <!-- Interests -->
                                                <div
                                                    class="border p-6 rounded-lg shadow text-center"
                                                >
                                                    <h1
                                                        class="text-lg font-semibold mb-2"
                                                    >
                                                        Interests
                                                    </h1>
                                                    <NumberFlowGroup
                                                        style="--number-flow-char-height: 0.85em"
                                                        class="flex items-center gap-4 font-semibold"
                                                    >
                                                        <NumberFlow
                                                            value={groupedData[
                                                                "Interests"
                                                            ]
                                                                ? groupedData[
                                                                      "Interests"
                                                                  ]
                                                                : 0}
                                                            format={{
                                                                style: "currency",
                                                                currency: "EUR",
                                                                maximumFractionDigits: 0,
                                                            }}
                                                            locales={"it-IT"}
                                                            class="text-3xl"
                                                        />
                                                        <NumberFlow
                                                            value={house_price ==
                                                                0 ||
                                                            !groupedData[
                                                                "Interests"
                                                            ]
                                                                ? 0
                                                                : groupedData[
                                                                      "Interests"
                                                                  ] /
                                                                  house_price}
                                                            format={{
                                                                style: "percent",
                                                                maximumFractionDigits: 2,
                                                                minimumFractionDigits: 2,
                                                                signDisplay:
                                                                    "always",
                                                            }}
                                                            class="text-lg transition-colors duration-300 text-red-500"
                                                        />
                                                    </NumberFlowGroup>
                                                </div>

                                                <!-- Mortgage Installment -->
                                                <div
                                                    class="border p-6 rounded-lg shadow text-center"
                                                >
                                                    <h1
                                                        class="text-lg font-semibold mb-2"
                                                    >
                                                        Mortgage Installment
                                                    </h1>
                                                    <NumberFlowGroup
                                                        style="--number-flow-char-height: 0.85em"
                                                        class="flex items-center gap-4 font-semibold"
                                                    >
                                                        <NumberFlow
                                                            value={mortgageInstallment}
                                                            format={{
                                                                style: "currency",
                                                                currency: "EUR",
                                                                maximumFractionDigits: 0,
                                                            }}
                                                            locales={"it-IT"}
                                                            class="text-3xl"
                                                        />
                                                    </NumberFlowGroup>
                                                </div>
                                            {/if}
                                        </div>
                                    </div>
                                </section>
                            {/if}
                            {#if selectedTab == "base"}{/if}
                            {#if selectedTab == "mortgage"}{/if}
                        </div>
                    {/if}
                </div>
            </main>
        </div>
    </div>
</div>

<style>
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
