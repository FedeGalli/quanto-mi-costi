<script lang="ts">
    import { onMount, tick } from "svelte";
    import Chart from "chart.js/auto";
    import { fade, slide } from "svelte/transition";
    import NumberFlow from "@number-flow/svelte";
    import Check from "../assets/Check.svelte";
    import NavBar from "../assets/NavBar.svelte";
    import DetailTile from "../assets/DetailTile.svelte";
    import TotalPriceTile from "../assets/TotalPriceTile.svelte";
    import SummaryDeltaPriceTile from "../assets/SummaryDeltaPriceTile.svelte";
    import DeltaPriceTile from "../assets/DeltaPriceTile.svelte";

    // Define types for the house cost breakdown
    type CostItem = {
        name: string;
        category: string;
        value: number;
        estimate: boolean;
        info: string;
    };

    let selectedTab = "summary";
    let prevSelectedTab = "";

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
    let interestsBarChart: Chart | null = null; // Store chart instance globally
    let dataByCategory: any = {};

    // Define category colors
    const categoryColors: Record<string, string> = {
        Tax: "rgba(255, 99, 132, 0.6)",
        Agency: "rgba(54, 162, 235, 0.6)",
        Notary: "rgba(201, 200, 115, 0.6)",
        Bank: "rgba(66, 245, 126, 0.6)",
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
            initializeBarChart();
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
            updateBarChart();
            updateInterestsBarChart();
        });
    }

    async function initializeBarChart() {
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

            // Step 4: Create the stacked bar chart
            barChartInstance = new Chart(ctx2D, {
                type: "bar",
                data: {
                    labels: ["Total Costs"],
                    datasets: [],
                },
                options: {
                    indexAxis: "x",
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        x: {
                            stacked: true,
                            grid: { display: false },
                            ticks: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                },
                            },
                        },
                        y: {
                            stacked: true,
                            grid: { display: false },
                            afterBuildTicks: function (axis) {
                                const ticks = [0];
                                axis.ticks = ticks.map((v) => ({ value: v }));
                            },
                            ticks: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                },
                            },
                        },
                    },
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
                            backgroundColor: ["rgba(255, 99, 132, 0.5)"],
                            borderColor: ["rgba(255, 99, 132, 1)"],
                            borderWidth: 1,
                        },
                        {
                            label: "Capital",
                            data: [],
                            backgroundColor: ["rgba(32, 128, 212, 0.5)"],
                            borderColor: ["rgba(32, 128, 212, 1)"],
                            borderWidth: 1,
                        },
                    ],
                },
                options: {
                    scales: {
                        y: {
                            stacked: true,
                            beginAtZero: true,
                            ticks: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                },
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
                        },
                    },
                },
            });
        }
    }

    function updateBarChart() {
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
        const max = Math.round(totalAmount);

        barChartInstance!.options!.scales!.y!.max = max;
        barChartInstance!.options!.scales!.y!.afterBuildTicks = function (
            axis,
        ) {
            const ticks = tickList;
            axis.ticks = ticks.map((v) => ({ value: v }));
        };
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

                interests = (residualCapital * taeg) / 12;
                interestsSum += interests;
                interestsVal += interests;

                capital = mortgageInstallment - interests;
                capitalSum += capital;
                residualCapital -= capital;

                if (i % 12 == 0) {
                    labels.push(
                        (new Date().getFullYear() + (i / 12 - 1)).toString(),
                    );
                    interestsData.push(interestsSum);
                    capitalData.push(capitalSum);

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
                backgroundColor: ["rgba(255, 99, 132, 0.5)"],
                borderColor: ["rgba(255, 99, 132, 1)"],
                borderWidth: 1,
            },
            {
                label: "Capital",
                data: capitalData,
                backgroundColor: ["rgba(32, 128, 212, 0.5)"],
                borderColor: ["rgba(32, 128, 212, 1)"],
                borderWidth: 1,
            },
        ];

        interestsBarChart.update(); // Smoothly update the chart
    }
</script>

<div class="h-screen overflow-hidden">
    <div class="flex h-full">
        <div class="w-86 bg-gray-700 text-white p-4 flex flex-col">
            <h1 class="text-xl font-bold mb-6">🏡 Real House Cost Simulator</h1>
            <!-- Inputs Section -->
            <section
                class="w-full max-w-3xl p-6 flex flex-col gap-6"
                transition:slide={{ duration: 500 }}
            >
                <div
                    class="flex flex-col gap-4 mb-4"
                    transition:fade={{ duration: 500 }}
                >
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
                                    min="0"
                                    max="50"
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
                            <button
                                on:click={() => {
                                    showCosts = true;
                                }}
                                class="btn"
                            >
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
                        <div transition:fade={{ duration: 500 }}>
                            <NavBar bind:selectedTab />
                            {#if selectedTab == "summary"}
                                <!-- Entire summary section -->
                                <section
                                    class="flex flex-col w-full gap-6"
                                    transition:slide={{ duration: 500 }}
                                >
                                    <!-- Main Chart + Cost Breakdown Side-by-Side -->
                                    <div
                                        class="flex flex-row gap-6 justify-center items-stretch"
                                        transition:fade={{ duration: 500 }}
                                    >
                                        <!-- Chart -->
                                        <div
                                            class="border p-6 rounded-lg shadow flex items-center justify-center flex-[1]"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <canvas
                                                id="barChart"
                                                class="w-full h-full"
                                            ></canvas>
                                        </div>

                                        <!-- Cost Breakdown Column -->
                                        <div
                                            class="flex flex-col flex-[2] max-w-lg w-full relative gap-6"
                                        >
                                            <!-- Total (Main) Price -->
                                            <TotalPriceTile
                                                number={totalAmount +
                                                    house_price}
                                                dif={totalAmount}
                                                name={"💸 Final Total Price"}
                                            />

                                            <!-- Subtiles Container with Vertical Line -->
                                            <div
                                                class="relative ml-15 space-y-6"
                                            >
                                                <!-- Main Vertical Line -->
                                                <div
                                                    class="absolute left-0 top-0 bottom-7 w-px bg-gray-50"
                                                ></div>

                                                <SummaryDeltaPriceTile
                                                    name={"🇮🇹 Taxes"}
                                                    number={groupedData["Tax"]}
                                                    showVal={house_price != 0}
                                                    delta={groupedData["Tax"] /
                                                        house_price}
                                                />
                                                <SummaryDeltaPriceTile
                                                    name={"💼 Notary"}
                                                    number={groupedData[
                                                        "Notary"
                                                    ]}
                                                    showVal={house_price != 0}
                                                    delta={groupedData[
                                                        "Notary"
                                                    ] / house_price}
                                                />

                                                {#if is_sold_by_agency}
                                                    <!-- Agency -->
                                                    <SummaryDeltaPriceTile
                                                        name={"🏘️ Agency"}
                                                        number={groupedData[
                                                            "Agency"
                                                        ]}
                                                        showVal={house_price >
                                                            0 &&
                                                            groupedData[
                                                                "Agency"
                                                            ] != null}
                                                        delta={groupedData[
                                                            "Agency"
                                                        ] / house_price}
                                                    />
                                                {/if}

                                                {#if is_using_mortgage}
                                                    <!-- Bank Costs -->
                                                    <SummaryDeltaPriceTile
                                                        name={"🏦 Bank"}
                                                        number={groupedData[
                                                            "Bank"
                                                        ]}
                                                        showVal={groupedData[
                                                            "Bank"
                                                        ] != null}
                                                        delta={groupedData[
                                                            "Bank"
                                                        ] / house_price}
                                                    />
                                                {/if}
                                            </div>
                                        </div>
                                    </div>
                                </section>
                            {/if}
                            {#if selectedTab == "base"}
                                <div
                                    class="flex flex-wrap gap-6"
                                    transition:slide={{ duration: 500 }}
                                >
                                    <div
                                        class="w-full md:w-[48%] mb-6"
                                        transition:slide={{ duration: 500 }}
                                    >
                                        <DetailTile
                                            data={dataByCategory["Tax"]}
                                            title={"🇮🇹 Taxes"}
                                        />
                                    </div>
                                    <div
                                        class="w-full md:w-[48%] mb-6"
                                        transition:slide={{ duration: 500 }}
                                    >
                                        <DetailTile
                                            data={dataByCategory["Notary"]}
                                            title={"💼 Notary"}
                                        />
                                    </div>
                                    {#if "Agency" in dataByCategory}
                                        <div
                                            class="w-full md:w-[48%] mb-6"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <DetailTile
                                                data={dataByCategory["Agency"]}
                                                title={"🏘️ Agency"}
                                            />
                                        </div>
                                    {/if}
                                    {#if "Bank" in dataByCategory}
                                        <div
                                            class="w-full md:w-[48%] mb-6"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <DetailTile
                                                data={dataByCategory["Bank"]}
                                                title={"🏦 Bank"}
                                            />
                                        </div>
                                    {/if}
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
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <div
                                                class="w-full md:w-[48%] mb-6"
                                                transition:slide={{
                                                    duration: 500,
                                                }}
                                            >
                                                <TotalPriceTile
                                                    number={mortgageInstallment}
                                                    name={"Monthly installment"}
                                                    dif={0}
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
                                                        0 && is_using_mortgage}
                                                />
                                            </div>
                                        </div>
                                        <!-- Chart -->
                                        <div
                                            class="border p-6 rounded-lg mt-6 shadow flex items-center justify-center"
                                            transition:slide={{ duration: 500 }}
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
