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
    import ColoredSummaryPrice from "../assets/ColoredSummaryPrice.svelte";

    let selectedTab = "summary";
    let prevSelectedTab = "";
    let showRateLimitPopup = false;
    let showErrorPopup = false;

    let house_price: number = 300000;
    let displayed_house_price: number = 300000;
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
    let mortgage_duration: number = 0;
    let mortgage_amount: number = 0;
    let displayed_mortgage_amount: number = 1000;
    let interestsVal: number = 0;
    let mortgageInstallment: number = 0;
    let agencyFee: number = 0;
    let showCosts: boolean = false;
    let showTooltip: boolean[] = [false, false, false, false];

    let totalAmount: number = 0;
    let doughnutChartInstance: Chart | null = null; // Store chart instance globally
    let interestsBarChartInstance: Chart | null = null; // Store chart instance globally
    let lineChartInstance: Chart | null = null;
    let dataByCategory: any = {};

    let wealth: number[] = [];
    let yearlySaving: number = 20000;
    let yearlyGrowthgRate: number = 3;

    // Define category colors
    const categoryColors: Record<string, string> = {
        Tax: "rgba(133, 81, 182, 0.7)",
        Agency: "rgba(249, 166, 0, 0.7)",
        Notary: "rgba(98, 182, 170, 0.7)",
        Bank: "rgba(76, 51, 141, 0.7)",
    };
    const categoryBorderColors: Record<string, string> = {
        Tax: "rgba(133, 81, 182, 1)",
        Agency: "rgba(249, 166, 0, 1)",
        Notary: "rgba(98, 182, 170, 1)",
        Bank: "rgba(76, 51, 141, 1)",
    };

    $: displayed_mortgage_amount = mortgage_amount;

    // Step 1: Group data by category and sum values
    let groupedData: Record<string, number> = {};

    // Optional: auto-hide after X seconds
    $: if (showRateLimitPopup || showErrorPopup) {
        setTimeout(() => {
            showRateLimitPopup = false;
            showErrorPopup = false;
        }, 4000); // Hide after 4 seconds
    }

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

    $: if (
        selectedTab == "mortgage_compare" &&
        prevSelectedTab != "mortgage_compare" &&
        showCosts
    ) {
        prevSelectedTab = selectedTab;
        tick().then(() => {
            initializeLineChart();
            updatePremiumData();
        });
    }

    $: if (selectedTab == "base") {
        prevSelectedTab = selectedTab;
    }

    $: if (
        (selectedTab == "mortgage" || selectedTab == "mortgage_compare") &&
        !is_using_mortgage
    ) {
        prevSelectedTab = selectedTab;
        selectedTab = "summary";
    }

    $: if (mortgage_duration < 0) {
        mortgage_duration = 0;
    }

    $: if (mortgage_duration > 50) {
        mortgage_duration = 50;
    }

    $: if (taeg < 0) {
        taeg = 0;
    }

    $: if (taeg > 20) {
        taeg = 20;
    }

    $: if (mortgage_amount < 0) {
        mortgage_amount = 1000;
    }

    $: if (house_price != null && mortgage_amount > house_price) {
        mortgage_amount = house_price;
    }
    $: if (house_price < 0) {
        house_price = 0;
    }

    $: if (house_price > 1500000) {
        house_price = 1500000;
    }

    $: if (yearlySaving > 200000) {
        yearlySaving = 200000;
    }
    $: if (yearlySaving < 0) {
        yearlySaving = 0;
    }
    $: if (yearlyGrowthgRate > 50) {
        yearlySaving = 50;
    }
    $: if (yearlyGrowthgRate < -10) {
        yearlySaving = -10;
    }

    onMount(() => {});

    function buildApiString(): string {
        let apiStringUrl: string =
            "https://quanto-mi-costi-backend.onrender.com/get_house_costs?house_price=" +
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
                (mortgage_duration != null ? mortgage_duration : 0) +
                "&mortgage_TAEG=" +
                (taeg != null ? taeg / 100 : 0);

        return apiStringUrl;
    }

    async function updateData() {
        let chartData: any = [];

        const apiStringUrl: string = buildApiString();
        let responseCode: string = "200";

        const apiCall = async () => {
            try {
                const response = await fetch(apiStringUrl);

                if (!response.ok || response.status == 429) {
                    throw new Error(`${response.status}`);
                }
                chartData = await response.json();
                chartData = chartData["data"];
            } catch (error: any) {
                responseCode = error.message;
            }
        };

        apiCall().then(() => {
            if (responseCode == "200") {
                dataByCategory = {};

                chartData.forEach((item: any) => {
                    if (item.category == "Installment") {
                        mortgageInstallment = item.value;
                    }
                    if (item.category == "House price") {
                        displayed_house_price = item.value;
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
            } else if (responseCode == "429") {
                showRateLimitPopup = true;
            } else {
                showErrorPopup = true;
            }
        });
    }

    async function updatePremiumData() {
        let chartData: any = [];

        const apiStringUrl: string = buildApiString();
        let responseCode: string = "200";

        const apiCall = async () => {
            try {
                const response = await fetch(apiStringUrl);

                if (!response.ok || response.status == 429) {
                    throw new Error(`${response.status}`);
                }
                chartData = await response.json();
                chartData = chartData["data"];
            } catch (error: any) {
                responseCode = error.message;
            }
        };

        apiCall().then(() => {
            if (responseCode == "200") {
                dataByCategory = {};

                chartData.forEach((item: any) => {
                    if (item.category == "Installment") {
                        mortgageInstallment = item.value;
                    }
                    if (item.category == "House price") {
                        displayed_house_price = item.value;
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
                updateLineChart();
            } else if (responseCode == "429") {
                showRateLimitPopup = true;
            } else {
                showErrorPopup = true;
            }
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
            if (doughnutChartInstance) {
                doughnutChartInstance.destroy();
            }

            const categories = Object.keys(groupedData);
            const values = Object.values(groupedData);

            // Step 4: Create the stacked bar chart
            doughnutChartInstance = new Chart(ctx2D, {
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
                                    "rgba(200, 200, 200, 1)",
                            ), // fallback color
                            borderColor: categories.map(
                                (cat) =>
                                    categoryBorderColors[cat] ||
                                    "rgba(200, 200, 200, 0.6)",
                            ),
                            borderWidth: 2,
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
                                    lineHeight: 1.25,
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

            if (interestsBarChartInstance) {
                interestsBarChartInstance.destroy();
            }

            interestsBarChartInstance = new Chart(ictx2D, {
                type: "bar",
                data: {
                    labels: [],
                    datasets: [
                        {
                            label: "Interests",
                            data: [],
                            backgroundColor: ["rgba(133, 81, 182, 0.7)"],
                            borderColor: ["rgba(133, 81, 182, 1)"],
                            borderWidth: 2,
                        },
                        {
                            label: "Capital",
                            data: [],
                            backgroundColor: ["rgba(249, 166, 0, 0.7)"],
                            borderColor: ["rgba(249, 166, 0, 1)"],
                            borderWidth: 2,
                        },
                    ],
                },
                options: {
                    responsive: true,
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
                                    lineHeight: 1.25,
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
                                display: true,
                            },
                            border: {
                                display: true, // ❌ remove x-axis border
                            },
                        },
                        x: {
                            stacked: true,
                            ticks: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                    lineHeight: 1.25,
                                },
                            },
                            grid: {
                                display: false,
                            },
                            border: {
                                display: false, // ❌ remove x-axis border
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
                                    lineHeight: 1.25,
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

    function simulateSavings(
        yearlyIncome: number[],
        installment: number,
        duration: number,
        taeg: number,
        yearlyGrowthgRate: number,
    ): number[] {
        let savings = 0;
        const savingsOverTime: number[] = [];
        let houseSaving = house_price - mortgage_amount;
        let totalMortgageAmount = mortgage_amount;

        yearlyGrowthgRate = yearlyGrowthgRate / 100;
        taeg = taeg / 100;

        for (let year = 0; year < yearlyIncome.length; year++) {
            let leftover = 0;
            if (year < duration) {
                leftover = yearlyIncome[year] - installment;
                houseSaving += installment - totalMortgageAmount * taeg;
                totalMortgageAmount -= installment - totalMortgageAmount * taeg;
            } else {
                leftover = yearlyIncome[year];
            }

            savings =
                savings + leftover >= 0
                    ? (savings + leftover) * (1 + yearlyGrowthgRate)
                    : savings + leftover;
            savingsOverTime.push(Math.round(savings + houseSaving));
        }

        return savingsOverTime;
    }

    function calculateYearlyInstallment(
        amount: number,
        duration: number,
        TAEG: number,
    ): number {
        if (amount === 0 || duration === 0 || TAEG === 0) {
            return 0;
        }

        const monthlyRate = TAEG / 12 / 100;
        const numberOfPayments = duration * 12;
        const numerator =
            amount * monthlyRate * Math.pow(1 + monthlyRate, numberOfPayments);
        const denominator = Math.pow(1 + monthlyRate, numberOfPayments) - 1;

        return (numerator / denominator) * 12;
    }

    async function initializeLineChart() {
        if (showCosts) {
            await tick(); // Wait for DOM update before selecting the canvas
            const ctx = document.getElementById(
                "lineChart",
            ) as HTMLCanvasElement | null;
            if (!ctx) return; // Ensure ctx exists before using it
            const ctx2D = ctx.getContext("2d");
            if (!ctx2D) return; // Ensure context is available

            // Step 3: Destroy previous chart instance
            if (lineChartInstance) {
                lineChartInstance.destroy();
            }

            if (mortgage_amount == 0 || taeg == 0) {
            }

            const durations: number[] = [10, 20, 30];
            const years = 30;

            // Simulated monthly income – replace this with your actual array if dynamic
            const yearlyIncome = Array.from(
                { length: years },
                (_, i) => yearlySaving,
            );

            // Example mortgage installments – adjust based on your input logic
            const installments: { [key: number]: number } = {
                [durations[0]]: calculateYearlyInstallment(
                    mortgage_amount,
                    durations[0],
                    taeg,
                ),
                [durations[1]]: calculateYearlyInstallment(
                    mortgage_amount,
                    durations[1],
                    taeg,
                ),
                [durations[2]]: calculateYearlyInstallment(
                    mortgage_amount,
                    durations[2],
                    taeg,
                ),
            };

            const datasets = durations.map((duration, idx) => {
                const data = simulateSavings(
                    yearlyIncome,
                    installments[duration],
                    duration,
                    taeg,
                    yearlyGrowthgRate,
                );
                return {
                    label: `${duration}-year mortgage`,
                    data,
                    borderColor: [
                        "rgba(98, 182, 170, 1)",
                        "rgba(133, 81, 182, 1)",
                        "rgba(249, 166, 0, 1)",
                    ][idx],
                    backgroundColor: [
                        "rgba(98, 182, 170, 0.7)",
                        "rgba(133, 81, 182, 0.7)",
                        "rgba(249, 166, 0, 0.7)",
                    ][idx],
                    pointRadius: 1, // <<< smaller dots
                    pointHoverRadius: 4, // <<< slightly bigger on hover
                    tension: 0.3,
                };
            });

            wealth[0] = datasets[0].data[datasets[0].data.length - 1];
            wealth[1] = datasets[1].data[datasets[1].data.length - 1];
            wealth[2] = datasets[2].data[datasets[2].data.length - 1];

            lineChartInstance = new Chart(ctx2D, {
                type: "line",
                data: {
                    labels: Array.from(
                        { length: years },
                        (_, i) => `Month ${i + 1}`,
                    ),
                    datasets,
                },
                options: {
                    responsive: true,
                    scales: {
                        y: {
                            title: {
                                display: true,
                                text: "Patrimonio Posseduto (€)",
                            },
                            ticks: {
                                callback: function (value) {
                                    var y: number = +value;
                                    if (y >= 1000000) {
                                        return (
                                            (y / 1000000)
                                                .toFixed(1)
                                                .replace(".0", "") + "M"
                                        );
                                    } else if (y >= 1000) {
                                        return (y / 1000).toFixed(0) + "k";
                                    }
                                    return y;
                                },
                            },
                        },
                        x: {
                            ticks: {
                                callback: function (value, index) {
                                    if ((index + 1) % 2 === 0) {
                                        return this.getLabelForValue(
                                            value as number,
                                        );
                                    }
                                    // Show every Nth tick if needed or return ''
                                    return "";
                                },
                            },
                        },
                    },
                    plugins: {
                        legend: {
                            display: true,
                            position: "bottom",
                            labels: {
                                color: "white",
                                font: {
                                    weight: "bold",
                                    lineHeight: 1.25,
                                },
                                padding: 20,
                            },
                        },
                    },
                },
            });
        }
    }

    function updateDoughnutChart() {
        if (!doughnutChartInstance) return; // Ensure the chart exists

        // Step 2: Convert grouped data into datasets
        const categories = Object.keys(groupedData);
        const values = Object.values(groupedData);

        doughnutChartInstance.data.labels = categories;

        doughnutChartInstance.data.datasets[0].data = values; // Update dataset (values and colors)
        doughnutChartInstance.data.datasets[0].backgroundColor = categories.map(
            (cat) => categoryColors[cat] || "rgba(200, 200, 200, 0.6)", // fallback
        );
        doughnutChartInstance.data.datasets[0].borderColor = categories.map(
            (cat) => categoryBorderColors[cat] || "rgba(200, 200, 200, 0.6)", // fallback
        );
        doughnutChartInstance.update(); // Smoothly update the chart
    }

    function updateInterestsBarChart() {
        if (!interestsBarChartInstance) return; // Ensure the chart exists

        interestsVal = 0;
        const labels: string[] = [];
        const interestsData: number[] = [];
        const capitalData: number[] = [];
        let residualCapital = mortgage_amount;
        let interestsSum = 0;
        let capitalSum = 0;

        if (is_using_mortgage) {
            for (let i = 1; i <= mortgage_duration * 12; i++) {
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

        interestsBarChartInstance.data.labels = labels; // Update dataset (values and colors)
        interestsBarChartInstance.data.datasets = [
            {
                label: "Interessi",
                data: interestsData,
                backgroundColor: ["rgba(249, 166, 0, 0.7)"],
                borderColor: ["rgba(249, 166, 0, 1)"],
                borderWidth: 2,
            },
            {
                label: "Capitale",
                data: capitalData,
                backgroundColor: ["rgba(133, 81, 182, 0.7)"],
                borderColor: ["rgba(133, 81, 182, 1)"],
                borderWidth: 2,
            },
        ];
        interestsBarChartInstance.options!.scales!.y!.max =
            mortgageInstallment * 12;
        interestsBarChartInstance.options!.scales!.y!.ticks!.callback =
            function (value) {
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
            };
        interestsBarChartInstance.update(); // Smoothly update the chart
    }

    function updateLineChart() {
        if (!lineChartInstance) return;

        const durations: number[] = [10, 20, 30];
        const years = 30;

        // Simulated monthly income – replace this with your actual array if dynamic
        const yearlyIncome = Array.from(
            { length: years },
            (_, i) => yearlySaving,
        );

        // Example mortgage installments – adjust based on your input logic
        const installments: { [key: number]: number } = {
            [durations[0]]: calculateYearlyInstallment(
                mortgage_amount,
                durations[0],
                taeg,
            ),
            [durations[1]]: calculateYearlyInstallment(
                mortgage_amount,
                durations[1],
                taeg,
            ),
            [durations[2]]: calculateYearlyInstallment(
                mortgage_amount,
                durations[2],
                taeg,
            ),
        };

        const updatedDatasets = durations.map((duration, idx) => {
            const data = simulateSavings(
                yearlyIncome,
                installments[duration],
                duration,
                taeg,
                yearlyGrowthgRate,
            );

            console.log(data);

            return {
                label: `Mutuo ${duration}-anni`,
                data,
                borderColor: [
                    "rgba(98, 182, 170, 1)",
                    "rgba(133, 81, 182, 1)",
                    "rgba(249, 166, 0, 1)",
                ][idx],
                backgroundColor: [
                    "rgba(98, 182, 170, 0.7)",
                    "rgba(133, 81, 182, 0.7)",
                    "rgba(249, 166, 0, 0.7)",
                ][idx],
                pointRadius: 2, // <<< smaller dots
                pointHoverRadius: 4, // <<< slightly bigger on hover
                tension: 0.3,
            };
        });

        wealth[0] = updatedDatasets[0].data[updatedDatasets[0].data.length - 1];
        wealth[1] = updatedDatasets[1].data[updatedDatasets[1].data.length - 1];
        wealth[2] = updatedDatasets[2].data[updatedDatasets[2].data.length - 1];

        lineChartInstance.data.labels = Array.from(
            { length: years },
            (_, i) => `${i + 1}`,
        );
        lineChartInstance.data.datasets = updatedDatasets;
        lineChartInstance.update();
    }
</script>

<div
    class="min-h-screen bg-gradient-to-b from-purple-400 to-[#1e1f25] flex items-center justify-center p-2 sm:p-6"
>
    {#if showRateLimitPopup || showErrorPopup}
        <div
            class="fixed top-4 left-1/2 transform -translate-x-1/2 bg-red-500 text-white px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg z-50 transition-all duration-300 animate-fadeInDown max-w-[90vw]"
            transition:slide={{ duration: 300 }}
        >
            <div class="flex items-center justify-between gap-2 sm:gap-4">
                <span class="text-sm sm:text-base"
                    >{showRateLimitPopup
                        ? "Troppe richieste. Riprova tra qualche secondo."
                        : "Qualcosa è andato storto. Riprova tra qualche secondo."}</span
                >
                <button
                    on:click={() => {
                        showRateLimitPopup = false;
                        showErrorPopup = false;
                    }}
                    class="text-white hover:text-red-200 font-bold text-lg leading-none min-w-[24px]"
                >
                    ×
                </button>
            </div>
        </div>
    {/if}

    <!-- Mobile-first responsive container -->
    <div class="w-full max-w-[1350px] mx-auto">
        <!-- Mobile: Stack vertically, Desktop: Side by side -->
        <div class="flex flex-col lg:flex-row gap-4 sm:gap-6 min-h-[80vh]">
            <!-- Left Panel: Form inputs -->
            <div
                class="w-full lg:basis-[40%] bg-[#1e1f25] text-white rounded-2xl shadow-lg p-4 sm:p-6 lg:p-8 overflow-auto"
            >
                <h1
                    class="text-2xl sm:text-3xl lg:text-4xl font-bold leading-tight mb-4"
                >
                    Scopri quanto costa,<br />
                    <span class="text-purple-400"> per davvero, </span><br />
                    la casa dei tuoi sogni. 🏡
                </h1>
                <p class="text-gray-400 text-xs sm:text-sm mt-2 mb-4">
                    Il prezzo finale mostrato indica il costo totale il giorno
                    del rogito. *Per spese variabili come il notaio/banca è
                    stata fatta una stima veritiera su tutte le voci di costo.
                </p>

                <div class="space-y-4 sm:space-y-6">
                    <!-- House Price Input -->
                    <div class="flex flex-col">
                        <h2
                            class="text-base sm:text-lg font-semibold leading-tight mb-2"
                        >
                            Costo casa:
                        </h2>
                        <div class="relative">
                            <input
                                type="number"
                                inputmode="decimal"
                                class="w-full border border-white rounded px-4 pr-12 py-3 sm:py-2 text-left font-medium text-white bg-transparent text-base sm:text-sm"
                                min="0"
                                max="1000000"
                                step="5000"
                                bind:value={house_price}
                                on:mouseup={showCosts &&
                                selectedTab != "mortgage"
                                    ? updateData
                                    : () => {}}
                            />
                            <!-- Separator Line -->
                            <div
                                class="absolute inset-y-0 right-10 w-px bg-white"
                            ></div>
                            <!-- Euro Symbol -->
                            <div
                                class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                            >
                                €
                            </div>
                        </div>
                    </div>

                    <!-- Checkboxes -->
                    <div id="isFirstHouse">
                        <Check
                            bind:label={is_fisrt_house_label}
                            bind:bind={is_fisrt_house}
                            bind:showTooltip
                            element={0}
                            on:change={showCosts && selectedTab != "mortgage"
                                ? updateData
                                : () => {}}
                        />
                    </div>

                    <div id="isSoldByBuilder">
                        <Check
                            bind:label={is_sold_by_builder_label}
                            bind:bind={is_sold_by_builder}
                            bind:showTooltip
                            element={1}
                            on:change={showCosts && selectedTab != "mortgage"
                                ? updateData
                                : () => {}}
                        />
                    </div>

                    <div id="isSoldByAgency">
                        <Check
                            bind:label={is_sold_by_agency_label}
                            bind:bind={is_sold_by_agency}
                            bind:showTooltip
                            element={2}
                            on:change={showCosts && selectedTab != "mortgage"
                                ? updateData
                                : () => {}}
                        />
                        {#if is_sold_by_agency}
                            <div
                                class="flex flex-col mt-3"
                                transition:slide={{ duration: 300 }}
                            >
                                <h3
                                    class="text-base font-semibold leading-tight mb-2"
                                >
                                    Provvigione:
                                </h3>
                                <div class="relative">
                                    <input
                                        type="number"
                                        class="w-full border border-white rounded px-4 pr-12 py-3 sm:py-2 text-left font-medium text-white bg-transparent text-base sm:text-sm"
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
                                        class="absolute inset-y-0 right-10 w-px bg-white"
                                    ></div>
                                    <!-- % Symbol -->
                                    <div
                                        class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                    >
                                        %
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
                                class="space-y-2"
                                transition:slide={{ duration: 300 }}
                            >
                                <!-- Mortgage Amount - Full Width -->
                                <div class="flex flex-col">
                                    <h3
                                        class="text-base font-semibold leading-tight mb-1"
                                    >
                                        Importo mutuo:
                                    </h3>
                                    <div class="relative">
                                        <input
                                            type="number"
                                            class="w-full border border-white rounded px-4 pr-12 py-3 sm:py-2 text-left font-medium text-white bg-transparent text-base sm:text-sm"
                                            min="0"
                                            max={house_price == null
                                                ? 0
                                                : house_price}
                                            step="5000"
                                            bind:value={mortgage_amount}
                                            on:change={showCosts
                                                ? () => {
                                                      updateData();
                                                      updatePremiumData();
                                                  }
                                                : () => {}}
                                        />
                                        <!-- Separator Line -->
                                        <div
                                            class="absolute inset-y-0 right-10 w-px bg-white"
                                        ></div>
                                        <!-- Euro Symbol -->
                                        <div
                                            class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                        >
                                            €
                                        </div>
                                    </div>
                                </div>

                                <!-- Duration and TAEG - Side by Side -->
                                <div
                                    class="grid grid-cols-1 sm:grid-cols-2 gap-4"
                                >
                                    <!-- Duration -->
                                    <div class="flex flex-col">
                                        <h3
                                            class="text-base font-semibold leading-tight mb-2"
                                        >
                                            Durata:
                                        </h3>
                                        <div class="relative">
                                            <input
                                                type="number"
                                                class="w-full border border-white rounded px-4 pr-12 py-3 sm:py-2 text-left font-medium text-white bg-transparent text-base sm:text-sm"
                                                min="0"
                                                bind:value={mortgage_duration}
                                                on:change={showCosts
                                                    ? updateData
                                                    : () => {}}
                                            />
                                            <!-- Separator Line -->
                                            <div
                                                class="absolute inset-y-0 right-10 w-px bg-white"
                                            ></div>
                                            <!-- A Symbol -->
                                            <div
                                                class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                            >
                                                A
                                            </div>
                                        </div>
                                    </div>

                                    <!-- TAEG -->
                                    <div class="flex flex-col">
                                        <h3
                                            class="text-base font-semibold leading-tight mb-2"
                                        >
                                            TAEG:
                                        </h3>
                                        <div class="relative">
                                            <input
                                                type="number"
                                                class="w-full border border-white rounded px-4 pr-12 py-3 sm:py-2 text-left font-medium text-white bg-transparent text-base sm:text-sm"
                                                min="0"
                                                step="0.1"
                                                bind:value={taeg}
                                                on:change={showCosts
                                                    ? () => {
                                                          updateData();
                                                          updatePremiumData();
                                                      }
                                                    : () => {}}
                                            />
                                            <!-- Separator Line -->
                                            <div
                                                class="absolute inset-y-0 right-10 w-px bg-white"
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
                            </div>
                        {/if}
                    </div>
                </div>
                {#if !showCosts}
                    <div
                        transition:slide={{ duration: 300 }}
                        class="overflow-hidden mt-6"
                    >
                        <div
                            class="flex items-center space-x-2 transition-all duration-300"
                        >
                            <CustomButton
                                title={"Quanto mi costi?"}
                                onClick={() => {
                                    for (
                                        let i = 0;
                                        i < showTooltip.length;
                                        i++
                                    ) {
                                        showTooltip[i] = false;
                                    }
                                    showCosts = true;
                                }}
                            ></CustomButton>
                        </div>
                    </div>
                {/if}
            </div>

            <!-- Right Panel: Results -->
            <div class="relative w-full lg:basis-[55%]">
                <!-- Right Content -->
                <div
                    class="bg-[#1e1f25] rounded-2xl shadow-lg p-4 sm:p-6 lg:p-8 overflow-auto w-full h-full min-h-[400px] lg:min-h-[600px] transition-all duration-600 ease-out-[cubic-bezier(0.22, 1, 0.36, 1)]"
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
                                    <!-- Main Chart + Cost Breakdown -->
                                    <div transition:fade={{ duration: 500 }}>
                                        <!-- Updated layout - horizontal on desktop, vertical on mobile -->
                                        <div
                                            class="flex flex-col lg:flex-row gap-6 sm:gap-8 justify-center items-center w-full h-full"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <!-- Chart - takes up space on the left with subtle enhancement -->
                                            <div
                                                class="flex justify-center w-full lg:w-1/2 max-w-full h-[200px] sm:h-[250px] lg:h-[300px] relative group"
                                            >
                                                <!-- Subtle background enhancement -->
                                                <div
                                                    class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent rounded-lg -z-10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                ></div>
                                                <canvas
                                                    id="barChart"
                                                    class="rounded-lg"
                                                ></canvas>
                                            </div>

                                            <!-- Total Price - takes up space on the right with hover effect -->
                                            <div
                                                class="w-full lg:w-1/2 flex justify-center group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg w-full"
                                                >
                                                    <!-- Subtle gradient overlay -->
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <TotalPriceTile
                                                        number={totalAmount +
                                                            displayed_house_price}
                                                        dif={totalAmount}
                                                        name={"💰 Spesa totale"}
                                                    />
                                                </div>
                                            </div>
                                        </div>

                                        <!-- Option 1: Simple 2x2 grid with enhanced hover effects -->
                                        <div
                                            class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4 sm:mt-6 w-full"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <!-- Tasse tile with hover effect -->
                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <!-- Subtle gradient overlay -->
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <SummaryDeltaPriceTile
                                                        name={"🇮🇹 Tasse"}
                                                        number={groupedData[
                                                            "Tax"
                                                        ]}
                                                        showVal={house_price !=
                                                            0}
                                                        delta={groupedData[
                                                            "Tax"
                                                        ] /
                                                            displayed_house_price}
                                                    />
                                                </div>
                                            </div>

                                            <!-- Notaio tile with hover effect -->
                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <!-- Subtle gradient overlay -->
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-teal-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <SummaryDeltaPriceTile
                                                        name={"📜 Notaio"}
                                                        number={groupedData[
                                                            "Notary"
                                                        ]}
                                                        showVal={house_price !=
                                                            0}
                                                        delta={groupedData[
                                                            "Notary"
                                                        ] /
                                                            displayed_house_price}
                                                    />
                                                </div>
                                            </div>

                                            <!-- Agenzia tile with hover effect (conditional) -->
                                            {#if is_sold_by_agency}
                                                <div
                                                    class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                                >
                                                    <div
                                                        class="relative overflow-hidden rounded-lg"
                                                    >
                                                        <!-- Subtle gradient overlay -->
                                                        <div
                                                            class="absolute inset-0 bg-gradient-to-br from-amber-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                        ></div>
                                                        <SummaryDeltaPriceTile
                                                            name={"🏘️ Agenzia"}
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
                                                            ] /
                                                                displayed_house_price}
                                                        />
                                                    </div>
                                                </div>
                                            {/if}

                                            <!-- Banca tile with hover effect (conditional) -->
                                            {#if is_using_mortgage}
                                                <div
                                                    class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                                >
                                                    <div
                                                        class="relative overflow-hidden rounded-lg"
                                                    >
                                                        <!-- Subtle gradient overlay -->
                                                        <div
                                                            class="absolute inset-0 bg-gradient-to-br from-indigo-600/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                        ></div>
                                                        <SummaryDeltaPriceTile
                                                            name={"🏦 Banca"}
                                                            number={groupedData[
                                                                "Bank"
                                                            ]}
                                                            showVal={groupedData[
                                                                "Bank"
                                                            ] != null}
                                                            delta={groupedData[
                                                                "Bank"
                                                            ] /
                                                                displayed_house_price}
                                                        />
                                                    </div>
                                                </div>
                                            {/if}
                                        </div>
                                    </div>
                                {/if}

                                {#if selectedTab == "base"}
                                    <div transition:fade={{ duration: 500 }}>
                                        <!-- Detail tiles - single column layout -->
                                        <div
                                            class="grid grid-cols-1 gap-4"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative rounded-lg before:absolute before:inset-0 before:bg-gradient-to-br before:from-purple-500/10 before:to-transparent before:opacity-0 group-hover:before:opacity-100 before:transition-opacity before:duration-300 before:pointer-events-none"
                                                >
                                                    <DetailTile
                                                        data={dataByCategory}
                                                        element={"Tax"}
                                                        title={"🇮🇹 Tasse"}
                                                    />
                                                </div>
                                            </div>

                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative rounded-lg before:absolute before:inset-0 before:bg-gradient-to-br before:from-teal-400/10 before:to-transparent before:opacity-0 group-hover:before:opacity-100 before:transition-opacity before:duration-300 before:pointer-events-none"
                                                >
                                                    <DetailTile
                                                        data={dataByCategory}
                                                        element={"Notary"}
                                                        title={"📜 Notaio"}
                                                    />
                                                </div>
                                            </div>

                                            {#if "Agency" in dataByCategory}
                                                <div
                                                    class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                                >
                                                    <div
                                                        class="relative rounded-lg before:absolute before:inset-0 before:bg-gradient-to-br before:from-amber-500/10 before:to-transparent before:opacity-0 group-hover:before:opacity-100 before:transition-opacity before:duration-300 before:pointer-events-none"
                                                    >
                                                        <DetailTile
                                                            data={dataByCategory}
                                                            element={"Agency"}
                                                            title={"🏘️ Agenzia"}
                                                        />
                                                    </div>
                                                </div>
                                            {/if}

                                            {#if "Bank" in dataByCategory}
                                                <div
                                                    class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                                >
                                                    <div
                                                        class="relative rounded-lg before:absolute before:inset-0 before:bg-gradient-to-br before:from-indigo-700/10 before:to-transparent before:opacity-0 group-hover:before:opacity-100 before:transition-opacity before:duration-300 before:pointer-events-none"
                                                    >
                                                        <DetailTile
                                                            data={dataByCategory}
                                                            element={"Bank"}
                                                            title={"🏦 Banca"}
                                                        />
                                                    </div>
                                                </div>
                                            {/if}
                                        </div>
                                    </div>
                                {/if}

                                {#if selectedTab == "mortgage"}
                                    <!-- Mortgage section -->
                                    <div transition:fade={{ duration: 500 }}>
                                        <!-- Summary tiles - responsive grid -->
                                        <div
                                            class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6 mb-6"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <SummaryDeltaPriceTile
                                                number={mortgageInstallment}
                                                name={"💸 Rata mutuo mensile"}
                                                showVal={mortgageInstallment !=
                                                    null}
                                                delta={null}
                                            />
                                            <DeltaPriceTile
                                                number={is_using_mortgage
                                                    ? interestsVal
                                                    : 0}
                                                name={"📈 Interessi in"}
                                                duration={mortgage_duration}
                                                delta={interestsVal /
                                                    displayed_mortgage_amount}
                                                showVal={mortgage_amount !=
                                                    null &&
                                                    displayed_mortgage_amount !=
                                                        0 &&
                                                    is_using_mortgage}
                                            />
                                        </div>

                                        <!-- Chart title -->
                                        <div
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <h2
                                                class="font-bold leading-tight text-lg sm:text-xl mb-4"
                                            >
                                                Cosa restituisco ogni anno?
                                            </h2>

                                            <!-- Chart - responsive sizing -->

                                            <div
                                                class="w-full h-[300px] sm:h-[350px] lg:h-[350px] relative"
                                            >
                                                <!-- Subtle background enhancement -->
                                                <div
                                                    class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent rounded-lg -z-10"
                                                ></div>
                                                <canvas
                                                    id="interestsBarChart"
                                                    class="absolute top-0 left-0 w-full h-full"
                                                ></canvas>
                                            </div>
                                        </div>
                                    </div>
                                {/if}

                                {#if selectedTab == "mortgage_compare"}
                                    <!-- Mortgage comparison section -->
                                    <div transition:fade={{ duration: 500 }}>
                                        <!-- Input section - responsive layout -->
                                        <div
                                            class="flex flex-col sm:flex-row sm:flex-wrap items-start sm:items-center gap-2 sm:gap-1 mb-6"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <div
                                                class="flex flex-col sm:flex-row sm:flex-wrap items-start sm:items-center gap-2 sm:gap-1"
                                            >
                                                <h3
                                                    class="font-bold text-white text-sm sm:text-base"
                                                >
                                                    Ad oggi risparmio/amo
                                                    all'anno:
                                                </h3>

                                                <!-- First input with € - enhanced -->
                                                <div
                                                    class="relative w-full sm:w-32 group"
                                                >
                                                    <input
                                                        type="number"
                                                        class="w-full border border-white rounded px-3 pr-10 py-2 text-left font-medium text-white bg-transparent transition-all duration-200 hover:border-white/80 focus:border-white focus:shadow-lg focus:shadow-white/20"
                                                        min="0"
                                                        step="2500"
                                                        bind:value={
                                                            yearlySaving
                                                        }
                                                        on:change={showCosts
                                                            ? updatePremiumData
                                                            : () => {}}
                                                    />
                                                    <div
                                                        class="absolute inset-y-0 right-9 w-px bg-white transition-colors duration-200 group-focus-within:bg-white/80"
                                                    ></div>
                                                    <div
                                                        class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                                    >
                                                        €
                                                    </div>
                                                </div>

                                                <h3
                                                    class="font-bold text-white text-sm sm:text-base"
                                                >
                                                    e fruttano il :
                                                </h3>

                                                <!-- Second input with % - enhanced -->
                                                <div
                                                    class="relative w-full sm:w-24 group"
                                                >
                                                    <input
                                                        type="number"
                                                        class="w-full border border-white rounded px-3 pr-10 py-2 text-left font-medium text-white bg-transparent transition-all duration-200 hover:border-white/80 focus:border-white focus:shadow-lg focus:shadow-white/20"
                                                        min="0"
                                                        step="0.1"
                                                        bind:value={
                                                            yearlyGrowthgRate
                                                        }
                                                        on:change={showCosts
                                                            ? updatePremiumData
                                                            : () => {}}
                                                    />
                                                    <div
                                                        class="absolute inset-y-0 right-9 w-px bg-white transition-colors duration-200 group-focus-within:bg-white/80"
                                                    ></div>
                                                    <div
                                                        class="absolute inset-y-0 right-3 flex items-center text-white font-semibold"
                                                    >
                                                        %
                                                    </div>
                                                </div>

                                                <h3
                                                    class="font-bold text-white text-sm sm:text-base"
                                                >
                                                    annuo.
                                                </h3>
                                            </div>

                                            <h2
                                                class="font-bold text-white leading-tight text-base sm:text-lg mt-4 sm:mt-2"
                                            >
                                                Qual è la durata del mutuo che
                                                più mi farà avere più patrimonio
                                                fra 30 anni?
                                            </h2>
                                        </div>

                                        <!-- Chart - responsive sizing with subtle enhancement -->
                                        <div
                                            class="w-full h-[300px] sm:h-[300px] lg:h-[350px] mb-6 relative"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <!-- Subtle background enhancement -->
                                            <div
                                                class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent rounded-lg -z-10"
                                            ></div>
                                            <canvas
                                                id="lineChart"
                                                class="w-full h-full rounded-lg"
                                            ></canvas>
                                        </div>

                                        <!-- Summary tiles - enhanced with subtle effects -->
                                        <div
                                            class="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-6"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <!-- Enhanced cards with subtle hover effects -->
                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <!-- Subtle gradient overlay -->
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-emerald-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth[0]}
                                                        name={"Mutuo 10 Anni"}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(98, 182, 170, 1)"}
                                                    />
                                                </div>
                                            </div>

                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <!-- Subtle gradient overlay -->
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth[1]}
                                                        name={"Mutuo 20 Anni"}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(133, 81, 182, 1)"}
                                                    />
                                                </div>
                                            </div>

                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <!-- Subtle gradient overlay -->
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-amber-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth[2]}
                                                        name={"Mutuo 30 Anni"}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(249, 166, 0, 1)"}
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                {/if}
                            </div>
                        {/if}
                    </div>
                </div>

                <!-- Tooltips (Overlay) -->
                <TooltipFirstHouse {showTooltip} index={0} />
                <TooltipVat {showTooltip} index={1} />
                <TooltipAgency {showTooltip} index={2} />
                <TooltipMortgage {showTooltip} index={3} />
            </div>
        </div>
    </div>
</div>

<style>
</style>
