<script>
    import { onMount } from "svelte";
    import { fade, slide } from "svelte/transition";

    // Input variables
    let yearlySaving = 50000;
    let yearlySavingRate = 20;
    let yearlyGrowthRate = 7;
    let mortgageTAEG = 3.5;
    let housePrice = 300000;

    // Chart variables
    let canvas;
    let chart;

    // Calculated data
    let wealthData = [];
    let chartLabels = [];

    // Generate years for 30-year projection
    function generateYears() {
        return Array.from({ length: 31 }, (_, i) => i);
    }

    // Calculate wealth for different mortgage scenarios
    function calculateWealth() {
        const years = generateYears();
        const annualSavings = yearlySaving * (yearlySavingRate / 100);
        const growthRate = yearlyGrowthRate / 100;
        const mortgageRate = mortgageTAEG / 100;

        const scenarios = [
            {
                name: "Solo Contanti",
                downPayment: 1.0,
                color: "rgba(34, 197, 94, 1)",
            },
            {
                name: "Mutuo 25%",
                downPayment: 0.75,
                color: "rgba(59, 130, 246, 1)",
            },
            {
                name: "Mutuo 50%",
                downPayment: 0.5,
                color: "rgba(147, 51, 234, 1)",
            },
            {
                name: "Mutuo 100%",
                downPayment: 0.0,
                color: "rgba(239, 68, 68, 1)",
            },
        ];

        wealthData = scenarios.map((scenario) => {
            const initialCash = housePrice * scenario.downPayment;
            const loanAmount = housePrice * (1 - scenario.downPayment);

            // Calculate monthly mortgage payment if there's a loan
            let monthlyPayment = 0;
            if (loanAmount > 0) {
                const monthlyRate = mortgageRate / 12;
                const numPayments = 30 * 12;
                monthlyPayment =
                    (loanAmount *
                        (monthlyRate *
                            Math.pow(1 + monthlyRate, numPayments))) /
                    (Math.pow(1 + monthlyRate, numPayments) - 1);
            }
            const annualMortgagePayment = monthlyPayment * 12;

            return {
                label: scenario.name,
                data: years.map((year) => {
                    if (year === 0) {
                        return housePrice - initialCash; // Initial wealth (house value minus cash paid)
                    }

                    // Available cash for investment after mortgage payments
                    const availableForInvestment =
                        annualSavings - annualMortgagePayment;

                    // Calculate accumulated investments
                    let investmentValue = 0;
                    for (let y = 1; y <= year; y++) {
                        const yearlyInvestment = Math.max(
                            0,
                            availableForInvestment,
                        );
                        investmentValue =
                            (investmentValue + yearlyInvestment) *
                            (1 + growthRate);
                    }

                    // House appreciation (assuming 2% annually)
                    const houseValue = housePrice * Math.pow(1.02, year);

                    // Remaining mortgage balance
                    let remainingMortgage = 0;
                    if (loanAmount > 0 && year < 30) {
                        const monthsElapsed = year * 12;
                        const monthlyRate = mortgageRate / 12;
                        const numPayments = 30 * 12;
                        remainingMortgage =
                            (loanAmount *
                                (Math.pow(1 + monthlyRate, numPayments) -
                                    Math.pow(1 + monthlyRate, monthsElapsed))) /
                            (Math.pow(1 + monthlyRate, numPayments) - 1);
                    }

                    // Total wealth = house value + investments - remaining mortgage
                    return houseValue + investmentValue - remainingMortgage;
                }),
                borderColor: scenario.color,
                backgroundColor: scenario.color.replace("1)", "0.1)"),
                fill: false,
                tension: 0.4,
                borderWidth: 3,
                pointRadius: 0,
                pointHoverRadius: 6,
            };
        });

        chartLabels = years;
        updateChart();
    }

    // Update chart
    function updateChart() {
        if (!chart || !wealthData.length) return;

        chart.data.labels = chartLabels;
        chart.data.datasets = wealthData;
        chart.update("none");
    }

    // Initialize chart
    function initChart() {
        if (typeof Chart === "undefined") {
            setTimeout(initChart, 100);
            return;
        }

        const ctx = canvas.getContext("2d");

        chart = new Chart(ctx, {
            type: "line",
            data: {
                labels: chartLabels,
                datasets: wealthData,
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true,
                        position: "top",
                        labels: {
                            color: "white",
                            font: {
                                size: 12,
                                weight: "bold",
                            },
                            usePointStyle: true,
                            pointStyle: "line",
                        },
                    },
                    tooltip: {
                        mode: "index",
                        intersect: false,
                        backgroundColor: "rgba(0, 0, 0, 0.8)",
                        titleColor: "white",
                        bodyColor: "white",
                        borderColor: "rgba(255, 255, 255, 0.2)",
                        borderWidth: 1,
                        callbacks: {
                            label: function (context) {
                                return (
                                    context.dataset.label +
                                    ": €" +
                                    new Intl.NumberFormat("it-IT").format(
                                        Math.round(context.parsed.y),
                                    )
                                );
                            },
                        },
                    },
                },
                scales: {
                    x: {
                        title: {
                            display: true,
                            text: "Anni",
                            color: "white",
                            font: {
                                size: 14,
                                weight: "bold",
                            },
                        },
                        grid: {
                            color: "rgba(255, 255, 255, 0.1)",
                        },
                        ticks: {
                            color: "white",
                            font: {
                                size: 11,
                            },
                        },
                    },
                    y: {
                        title: {
                            display: true,
                            text: "Patrimonio (€)",
                            color: "white",
                            font: {
                                size: 14,
                                weight: "bold",
                            },
                        },
                        grid: {
                            color: "rgba(255, 255, 255, 0.1)",
                        },
                        ticks: {
                            color: "white",
                            font: {
                                size: 11,
                            },
                            callback: function (value) {
                                return (
                                    "€" +
                                    new Intl.NumberFormat("it-IT", {
                                        notation: "compact",
                                        maximumFractionDigits: 1,
                                    }).format(value)
                                );
                            },
                        },
                    },
                },
                interaction: {
                    mode: "index",
                    intersect: false,
                },
                elements: {
                    point: {
                        hoverBackgroundColor: "white",
                    },
                },
            },
        });
    }

    // Load Chart.js and initialize
    onMount(() => {
        // Load Chart.js from CDN
        const script = document.createElement("script");
        script.src =
            "https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.9.1/chart.min.js";
        script.onload = () => {
            calculateWealth();
            initChart();
        };
        document.head.appendChild(script);

        return () => {
            if (chart) {
                chart.destroy();
            }
        };
    });

    // Recalculate when inputs change
    $: if (
        yearlySaving ||
        yearlySavingRate ||
        yearlyGrowthRate ||
        mortgageTAEG ||
        housePrice
    ) {
        calculateWealth();
    }
</script>

<div
    class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 p-4"
>
    <div class="max-w-6xl mx-auto">
        <!-- Title -->
        <h1 class="text-3xl font-bold text-white text-center mb-8">
            Confronto Patrimonio: Contanti vs Mutuo
        </h1>

        <!-- Input Section -->
        <div
            class="bg-white/10 backdrop-blur-sm rounded-xl p-6 mb-8"
            transition:fade={{ duration: 500 }}
        >
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <!-- House Price -->
                <div class="flex flex-col items-center gap-2">
                    <h3 class="font-bold text-white text-sm text-center">
                        Prezzo Casa
                    </h3>
                    <div class="relative">
                        <input
                            type="number"
                            class="w-32 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-purple-400 focus:outline-none"
                            min="0"
                            step="10000"
                            bind:value={housePrice}
                        />
                        <span
                            class="absolute right-3 top-2 text-white font-semibold"
                            >€</span
                        >
                    </div>
                </div>

                <!-- Yearly Income -->
                <div class="flex flex-col items-center gap-2">
                    <h3 class="font-bold text-white text-sm text-center">
                        Guadagno Annuale
                    </h3>
                    <div class="relative">
                        <input
                            type="number"
                            class="w-32 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-purple-400 focus:outline-none"
                            min="0"
                            step="2500"
                            bind:value={yearlySaving}
                        />
                        <span
                            class="absolute right-3 top-2 text-white font-semibold"
                            >€</span
                        >
                    </div>
                </div>

                <!-- Savings Rate -->
                <div class="flex flex-col items-center gap-2">
                    <h3 class="font-bold text-white text-sm text-center">
                        Tasso di Risparmio
                    </h3>
                    <div class="relative">
                        <input
                            type="number"
                            class="w-24 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-purple-400 focus:outline-none"
                            min="0"
                            max="100"
                            step="1"
                            bind:value={yearlySavingRate}
                        />
                        <span
                            class="absolute right-3 top-2 text-white font-semibold"
                            >%</span
                        >
                    </div>
                </div>

                <!-- Investment Growth Rate -->
                <div class="flex flex-col items-center gap-2">
                    <h3 class="font-bold text-white text-sm text-center">
                        Rendimento Investimenti
                    </h3>
                    <div class="relative">
                        <input
                            type="number"
                            class="w-24 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-purple-400 focus:outline-none"
                            min="0"
                            step="0.1"
                            bind:value={yearlyGrowthRate}
                        />
                        <span
                            class="absolute right-3 top-2 text-white font-semibold"
                            >%</span
                        >
                    </div>
                </div>

                <!-- Mortgage TAEG -->
                <div class="flex flex-col items-center gap-2">
                    <h3 class="font-bold text-white text-sm text-center">
                        TAEG Mutuo
                    </h3>
                    <div class="relative">
                        <input
                            type="number"
                            class="w-24 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-purple-400 focus:outline-none"
                            min="0"
                            step="0.1"
                            bind:value={mortgageTAEG}
                        />
                        <span
                            class="absolute right-3 top-2 text-white font-semibold"
                            >%</span
                        >
                    </div>
                </div>
            </div>
        </div>

        <!-- Chart Section -->
        <div
            class="bg-white/10 backdrop-blur-sm rounded-xl p-6"
            transition:slide={{ duration: 500 }}
        >
            <h2 class="font-bold text-white text-xl text-center mb-6">
                Evoluzione del Patrimonio in 30 Anni
            </h2>

            <div class="w-full h-[400px] mb-6 relative">
                <div
                    class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent rounded-lg -z-10"
                ></div>
                <canvas bind:this={canvas} class="w-full h-full rounded-lg"
                ></canvas>
            </div>

            <!-- Summary Cards -->
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                {#each wealthData as scenario, index}
                    <div
                        class="bg-white/5 backdrop-blur-sm rounded-lg p-4 text-center border border-white/10 hover:border-white/30 transition-colors"
                    >
                        <div
                            class="w-4 h-4 rounded-full mx-auto mb-2"
                            style="background-color: {scenario.borderColor}"
                        ></div>
                        <h3 class="text-white font-semibold text-sm mb-1">
                            {scenario.label}
                        </h3>
                        <p class="text-white text-lg font-bold">
                            €{new Intl.NumberFormat("it-IT").format(
                                Math.round(
                                    scenario.data[scenario.data.length - 1],
                                ),
                            )}
                        </p>
                        <p class="text-white/70 text-xs">dopo 30 anni</p>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>

<style>
    input::-webkit-outer-spin-button,
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }

    input[type="number"] {
        -moz-appearance: textfield;
    }
</style>
