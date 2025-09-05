<script lang="ts">
    import { onMount, tick, onDestroy } from "svelte";
    import { push } from "svelte-spa-router";
    import Chart from "chart.js/auto";
    import { fade, fly, slide } from "svelte/transition";
    import Check from "../assets/Check.svelte";
    import NavBar from "../assets/NavBar.svelte";
    import DetailTile from "../assets/DetailTile.svelte";
    import TotalPriceTile from "../assets/TotalPriceTile.svelte";
    import SummaryDeltaPriceTile from "../assets/SummaryDeltaPriceTile.svelte";
    import DeltaPriceTile from "../assets/DeltaPriceTile.svelte";
    import TooltipFirstHouse from "../tooltip/TooltipFirstHouse.svelte";
    import TooltipVat from "../tooltip/TooltipVAT.svelte";
    import TooltipAgency from "../tooltip/TooltipAgency.svelte";
    import TooltipMortgage from "../tooltip/TooltipMortgage.svelte";
    import CustomButton from "../assets/CustomButton.svelte";
    import ColoredSummaryPrice from "../assets/ColoredSummaryPrice.svelte";
    import { onAuthStateChanged } from "firebase/auth";
    import { auth, db } from "./auth/credentials";
    import { doc, setDoc, getDoc } from "firebase/firestore";
    import {
        user,
        isAuthenticated,
        isLoading,
        initAuthStore,
        logout,
    } from "./auth/auth-store";
    import Prices from "./Prices.svelte";

    let selectedTab = "summary";

    let showRateLimitPopup = false;
    let showErrorPopup = false;
    let forceChartUpdate = false;
    let showUserMenu = false;

    // Handle logout
    async function handleLogout() {
        try {
            await logout();
            showUserMenu = false;
        } catch (error) {
            console.error("Logout failed:", error);
        }
    }

    function toggleUserMenu() {
        showUserMenu = !showUserMenu;
    }

    function closeUserMenu(event: any) {
        if (!event.target.closest(".user-menu-container")) {
            showUserMenu = false;
        }
    }

    onMount(() => {
        // Add click listener for closing user menu
        document.addEventListener("click", closeUserMenu);
    });

    onDestroy(() => {
        // Remove click listener
        document.removeEventListener("click", closeUserMenu);
    });

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
    let showMortgageParams = false; // Toggle state for mortgage parameters
    let mortgage_percentage: number[] = [0, 25, 50, 80];
    let mortgage_durations: number[] = [10, 20, 30];
    let cashVsMortgageData: any = [
        { percentage: 0, values: [0] },
        { percentage: 25, values: [0] },
        { percentage: 50, values: [0] },
        { percentage: 80, values: [0] },
    ];
    let mortgageCompareData: any = [
        { duration: 10, values: [0], valid: false, installment: 0 },
        { duration: 20, values: [0], valid: false, installment: 0 },
        { duration: 30, values: [0], valid: false, installment: 0 },
    ];

    let totalAmount: number = 0;
    let doughnutChartInstance: Chart | null = null; // Store chart instance globally
    let interestsBarChartInstance: Chart | null = null; // Store chart instance globally
    let lineChartInstance: Chart | null = null;
    let cashVsMortgageChartInstance: Chart | null = null;
    let dataByCategory: any = {};

    let wealth: number[] = [];
    let wealth_cash_vs_mortgage: number[] = [];
    let wealth_cash_vs_mortgage_installments: number[] = [];
    let yearlySaving: number = 20000;
    let yearlyGrowthgRate: number = 3;
    let yearlySavingRate: number = 30;

    let groupedData: Record<string, number> = {};

    // Define category colors
    const categoryColors: Record<string, string> = {
        Tasse: "rgba(133, 81, 182, 0.7)",
        Agenzia: "rgba(249, 166, 0, 0.7)",
        Notaio: "rgba(98, 182, 170, 0.7)",
        Banca: "rgba(76, 51, 141, 0.7)",
    };
    const categoryBorderColors: Record<string, string> = {
        Tasse: "rgba(133, 81, 182, 1)",
        Agenzia: "rgba(249, 166, 0, 1)",
        Notaio: "rgba(98, 182, 170, 1)",
        Banca: "rgba(76, 51, 141, 1)",
    };

    $: displayed_mortgage_amount = mortgage_amount;

    // Optional: auto-hide after X seconds
    $: if (showRateLimitPopup || showErrorPopup) {
        setTimeout(() => {
            showRateLimitPopup = false;
            showErrorPopup = false;
        }, 4000); // Hide after 4 seconds
    }

    //updating the chart when switching nav
    $: if (
        (selectedTab == "summary" && showCosts) ||
        (forceChartUpdate && showCosts && selectedTab == "summary")
    ) {
        forceChartUpdate = false;
        tick().then(() => {
            initializeDoughnutChart();
            updateData();
        });
    }

    $: if (
        (selectedTab == "mortgage" && showCosts) ||
        (forceChartUpdate && showCosts)
    ) {
        forceChartUpdate = false;
        tick().then(() => {
            initializeInterestsBarChart();
            updateData();
        });
    }

    $: if (
        (selectedTab == "mortgage_compare" && showCosts) ||
        (forceChartUpdate && showCosts && selectedTab == "mortgage_compare")
    ) {
        tick().then(() => {
            forceChartUpdate = false;
            initializeLineChart();
            updateMortgageCompare();
        });
    }

    $: if (
        (selectedTab == "cash_vs_mortgage" && showCosts) ||
        (forceChartUpdate && showCosts && selectedTab == "cash_vs_mortgage")
    ) {
        forceChartUpdate = false;
        tick().then(() => {
            initializeCashVsMortgageChart(); // to modify
            updateCashVsMortgage();
        });
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
    $: if (yearlySavingRate > 100) {
        yearlySavingRate = 100;
    }
    $: if (yearlySavingRate < 0) {
        yearlySavingRate = 0;
    }
    $: if (yearlyGrowthgRate > 50) {
        yearlyGrowthgRate = 50;
    }
    $: if (yearlyGrowthgRate < -10) {
        yearlyGrowthgRate = -10;
    }

    function saveStateToLocalStorage() {
        const state = {
            selectedTab,
            showRateLimitPopup,
            showErrorPopup,
            house_price,
            displayed_house_price,
            is_fisrt_house,
            is_sold_by_builder,
            is_sold_by_agency,
            is_using_mortgage,
            taeg,
            mortgage_duration,
            mortgage_amount,
            displayed_mortgage_amount,
            interestsVal,
            mortgageInstallment,
            agencyFee,
            showCosts,
            showTooltip,
            showMortgageParams,
            mortgage_percentage,
            mortgage_durations,
            cashVsMortgageData,
            mortgageCompareData,
            totalAmount,
            dataByCategory,
            wealth,
            wealth_cash_vs_mortgage,
            wealth_cash_vs_mortgage_installments,
            yearlySaving,
            yearlyGrowthgRate,
            yearlySavingRate,
            groupedData,
        };

        localStorage.setItem("appState", JSON.stringify(state));
    }
    function loadStateFromLocalStorage() {
        const savedState = localStorage.getItem("appState");
        if (savedState) {
            try {
                const state = JSON.parse(savedState);

                // Restore all variables
                selectedTab = state.selectedTab ?? "summary";
                showRateLimitPopup = state.showRateLimitPopup ?? false;
                showErrorPopup = state.showErrorPopup ?? false;
                house_price = state.house_price ?? 300000;
                displayed_house_price = state.displayed_house_price ?? 300000;
                is_fisrt_house = state.is_fisrt_house ?? false;
                is_sold_by_builder = state.is_sold_by_builder ?? false;
                is_sold_by_agency = state.is_sold_by_agency ?? false;
                is_using_mortgage = state.is_using_mortgage ?? false;
                taeg = state.taeg ?? 0;
                mortgage_duration = state.mortgage_duration ?? 0;
                mortgage_amount = state.mortgage_amount ?? 0;
                displayed_mortgage_amount =
                    state.displayed_mortgage_amount ?? 1000;
                interestsVal = state.interestsVal ?? 0;
                mortgageInstallment = state.mortgageInstallment ?? 0;
                agencyFee = state.agencyFee ?? 0;
                showCosts = state.showCosts ?? false;
                showTooltip = state.showTooltip ?? [false, false, false, false];
                showMortgageParams = state.showMortgageParams ?? false;
                mortgage_percentage = state.mortgage_percentage ?? [
                    0, 25, 50, 80,
                ];
                mortgage_durations = state.mortgage_durations ?? [10, 20, 30];
                cashVsMortgageData = state.cashVsMortgageData ?? [
                    { percentage: 0, values: [0] },
                    { percentage: 25, values: [0] },
                    { percentage: 50, values: [0] },
                    { percentage: 80, values: [0] },
                ];
                mortgageCompareData = state.mortgageCompareData ?? [
                    { duration: 10, values: [0], valid: false, installment: 0 },
                    { duration: 20, values: [0], valid: false, installment: 0 },
                    { duration: 30, values: [0], valid: false, installment: 0 },
                ];
                totalAmount = state.totalAmount ?? 0;
                dataByCategory = state.dataByCategory ?? {};
                wealth = state.wealth ?? [];
                wealth_cash_vs_mortgage = state.wealth_cash_vs_mortgage ?? [];
                wealth_cash_vs_mortgage_installments =
                    state.wealth_cash_vs_mortgage_installments ?? [];
                yearlySaving = state.yearlySaving ?? 20000;
                yearlyGrowthgRate = state.yearlyGrowthgRate ?? 3;
                yearlySavingRate = state.yearlySavingRate ?? 30;
                groupedData = state.groupedData ?? {};

                forceChartUpdate = true;
            } catch (error) {
                console.error("Error loading state from localStorage:", error);
            }
        }
    }
    function clearStateFromLocalStorage() {
        localStorage.removeItem("appState");
    }
    async function getUserData(firebaseUser: any) {
        const userData: any = {
            uid: firebaseUser.uid,
            email: firebaseUser.email,
            displayName: firebaseUser.displayName,
            photoURL: firebaseUser.photoURL,
            emailVerified: firebaseUser.emailVerified,
            firstName: firebaseUser.displayName?.split(" ")[0] || "",
            lastName:
                firebaseUser.displayName?.split(" ").slice(1).join(" ") || "",
            pro: await getUserPro(firebaseUser.uid),
        };
        return userData;
    }
    async function getUserPro(uid: string) {
        try {
            const userDocRef = doc(db, "users", uid);

            // Check if document already exists
            const userDoc = await getDoc(userDocRef);
            const userData = userDoc.data();

            if (userData) {
                return userData.is_pro;
            } else {
                return false;
            }
        } catch (error) {
            console.error("Error getting pro info:", error);
            // Don't throw error to avoid disrupting the login flow
        }
    }
    onMount(() => {
        initAuthStore();
        loadStateFromLocalStorage();
        clearStateFromLocalStorage();

        // Add this auth listener to your homepage too
        const unsubscribe = onAuthStateChanged(auth, (firebaseUser) => {
            if (firebaseUser && firebaseUser.emailVerified) {
                // User is signed in and email is verified
                getUserData(firebaseUser)
                    .then((userData) => {
                        user.set(userData);
                        isAuthenticated.set(true);
                    })
                    .catch((error) => {
                        console.error("Error initializing user:", error);
                    });
            }
        });

        return () => unsubscribe();
    });
    function goToLogIn() {
        saveStateToLocalStorage();
        push("/signin");
    }

    function goToPro() {
        saveStateToLocalStorage();
        push("/getpro");
    }
    function buildCostApiString(): string {
        let apiStringUrl: string =
            //"https://quanto-mi-costi-backend.onrender.com/get_house_costs?house_price=" +
            "http://localhost:8080/get_house_costs?house_price=" +
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
    function buildCashVsMortgageApiString(): string {
        let apiStringUrl: string =
            //"https://quanto-mi-costi-backend.onrender.com/get_house_costs?house_price=" +
            "http://localhost:8080/get_cash_vs_mortgage?mortgage_percentage=" +
            (mortgage_percentage != null ? mortgage_percentage.join(",") : "") +
            "&house_price=" +
            (house_price != null ? house_price : 0) +
            "&yearly_saving=" +
            (yearlySaving != null ? yearlySaving : 0) +
            "&yearly_growth_rate=" +
            (yearlyGrowthgRate != null ? yearlyGrowthgRate / 100 : 0) +
            "&yearly_saving_rate=" +
            (yearlySavingRate != null ? yearlySavingRate / 100 : 0) +
            "&mortgage_duration=" +
            (mortgage_duration != null ? mortgage_duration : 0) +
            "&mortgage_TAEG=" +
            (taeg != null ? taeg / 100 : 0);
        return apiStringUrl;
    }
    function buildMortgageCompareApiString(): string {
        let apiStringUrl: string =
            //"https://quanto-mi-costi-backend.onrender.com/get_house_costs?house_price=" +
            "http://localhost:8080/get_mortgage_compare?house_price=" +
            (house_price != null ? house_price : 0) +
            "&yearly_saving=" +
            (yearlySaving != null ? yearlySaving : 0) +
            "&yearly_growth_rate=" +
            (yearlyGrowthgRate != null ? yearlyGrowthgRate / 100 : 0) +
            "&yearly_saving_rate=" +
            (yearlySavingRate != null ? yearlySavingRate / 100 : 0) +
            "&mortgage_amount=" +
            (mortgage_amount != null ? mortgage_amount : 0) +
            "&mortgage_TAEG=" +
            (taeg != null ? taeg / 100 : 0) +
            "&durations=" +
            (mortgage_durations != null ? mortgage_durations.join(",") : "");
        return apiStringUrl;
    }
    function itTranslation(chartData: Record<string, number>) {
        let categories = Object.keys(chartData);

        for (let i = 0; i < categories.length; i++) {
            if (categories[i] == "Agency") {
                categories[i] = "Agenzia";
            } else if (categories[i] == "Notary") {
                categories[i] = "Notaio";
            } else if (categories[i] == "Bank") categories[i] = "Banca";
            else categories[i] = "Tasse";
        }

        return categories;
    }
    async function updateData() {
        let chartData: any = [];

        const apiStringUrl: string = buildCostApiString();
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
    async function updateMortgageCompare() {
        let responseCode: string = "200";

        const mortgageCompareApiCall = async () => {
            try {
                const response = await fetch(buildMortgageCompareApiString());

                if (!response.ok || response.status == 429) {
                    throw new Error(`${response.status}`);
                }
                mortgageCompareData = await response.json();
                mortgageCompareData = mortgageCompareData["data"];
            } catch (error: any) {
                responseCode = error.message;
            }
        };
        mortgageCompareApiCall().then(() => {
            if (responseCode == "200") {
                updateLineChart();
            } else if (responseCode == "429") {
                showRateLimitPopup = true;
            } else {
                showErrorPopup = true;
            }
        });
    }
    async function updateCashVsMortgage() {
        let responseCode: string = "200";

        const cashVsMortgageApiCall = async () => {
            try {
                const response = await fetch(buildCashVsMortgageApiString());

                if (!response.ok || response.status == 429) {
                    throw new Error(`${response.status}`);
                }
                cashVsMortgageData = await response.json();
                cashVsMortgageData = cashVsMortgageData["data"];
            } catch (error: any) {
                responseCode = error.message;
            }
        };
        cashVsMortgageApiCall().then(() => {
            if (responseCode == "200") {
                updateCashVsMortgageChart();
            } else if (responseCode == "429") {
                showRateLimitPopup = true;
            } else {
                showErrorPopup = true;
            }
        });
    }
    async function initializeDoughnutChart() {
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

        const categories = itTranslation(groupedData);
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
                                categoryColors[cat] || "rgba(200, 200, 200, 1)",
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
                                display: false,
                            },
                            border: {
                                display: false, // ❌ remove x-axis border
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

            const durations: number[] = [10, 20, 30];
            const years = 30 + 1;

            const datasets = durations.map((duration, idx) => {
                const data = mortgageCompareData[idx]["values"];
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
                        (_, i) => `Month ${i}`,
                    ),
                    datasets,
                },
                options: {
                    responsive: true,
                    interaction: {
                        intersect: false,
                        mode: "index",
                    },
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
                                    if (index % 2 === 0) {
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
                        tooltip: {
                            callbacks: {
                                title: function (context) {
                                    // Get the index from the first tooltip item
                                    const dataIndex = context[0].dataIndex;
                                    // Calculate the year (assuming 12 months per year)
                                    const year = Math.floor(dataIndex);

                                    // You can choose one of these formats:
                                    // Option 1: Just show the year
                                    return `Anno ${year}`;

                                    // Option 2: Show year and month
                                    // return `Anno ${year}, Mese ${month}`;

                                    // Option 3: Show year and month in different format
                                    // return `Anno ${year} - Mese ${month}`;
                                },
                            },
                        },
                    },
                },
            });
        }
    }
    async function initializeCashVsMortgageChart() {
        if (showCosts) {
            await tick(); // Wait for DOM update before selecting the canvas
            const ctx = document.getElementById(
                "cashVsMortgageChart",
            ) as HTMLCanvasElement | null;
            if (!ctx) return; // Ensure ctx exists before using it
            const ctx2D = ctx.getContext("2d");
            if (!ctx2D) return; // Ensure context is available

            // Step 3: Destroy previous chart instance
            if (cashVsMortgageChartInstance) {
                cashVsMortgageChartInstance.destroy();
            }

            const mortgage_percentage: number[] = [0, 25, 50, 80];
            const years = 30 + 1;

            let i = 0;
            const datasets = mortgage_percentage.map(
                (mortgage_percentage, idx) => {
                    const data = cashVsMortgageData[idx]["values"];
                    wealth_cash_vs_mortgage_installments[i] =
                        cashVsMortgageData[idx]["installment"];
                    i += 1;
                    return {
                        label: `${mortgage_percentage}% mutuo`,
                        data,
                        borderColor: [
                            "rgba(178, 178, 178, 1)",
                            "rgba(98, 182, 170, 1)",
                            "rgba(133, 81, 182, 1)",
                            "rgba(249, 166, 0, 1)",
                        ][idx],
                        backgroundColor: [
                            "rgba(178, 178, 178, 0.7)",
                            "rgba(98, 182, 170, 0.7)",
                            "rgba(133, 81, 182, 0.7)",
                            "rgba(249, 166, 0, 0.7)",
                        ][idx],
                        pointRadius: 1, // <<< smaller dots
                        pointHoverRadius: 4, // <<< slightly bigger on hover
                        tension: 0.3,
                    };
                },
            );

            wealth_cash_vs_mortgage[0] =
                datasets[0].data[datasets[0].data.length - 1];
            wealth_cash_vs_mortgage[1] =
                datasets[1].data[datasets[1].data.length - 1];
            wealth_cash_vs_mortgage[2] =
                datasets[2].data[datasets[2].data.length - 1];
            wealth_cash_vs_mortgage[3] =
                datasets[3].data[datasets[3].data.length - 1];

            cashVsMortgageChartInstance = new Chart(ctx2D, {
                type: "line",
                data: {
                    labels: Array.from(
                        { length: years },
                        (_, i) => `Month ${i}`,
                    ),
                    datasets,
                },
                options: {
                    responsive: true,
                    interaction: {
                        intersect: false,
                        mode: "index",
                    },
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
                                    if (index % 2 === 0) {
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
                        tooltip: {
                            callbacks: {
                                title: function (context) {
                                    // Get the index from the first tooltip item
                                    const dataIndex = context[0].dataIndex;
                                    // Calculate the year (assuming 12 months per year)
                                    const year = Math.floor(dataIndex);

                                    // You can choose one of these formats:
                                    // Option 1: Just show the year
                                    return `Anno ${year}`;

                                    // Option 2: Show year and month
                                    // return `Anno ${year}, Mese ${month}`;

                                    // Option 3: Show year and month in different format
                                    // return `Anno ${year} - Mese ${month}`;
                                },
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
        const categories = itTranslation(groupedData);
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
        const years = 30 + 1;

        const updatedDatasets = durations.map((duration, idx) => {
            const data = mortgageCompareData[idx]["values"];
            return {
                label: `Mutuo ${duration} anni`,
                data,
                borderColor: [
                    mortgageCompareData[idx]["valid"]
                        ? "rgba(98, 182, 170, 1)"
                        : "rgba(98, 182, 170, 0.2)",
                    mortgageCompareData[idx]["valid"]
                        ? "rgba(133, 81, 182, 1)"
                        : "rgba(133, 81, 182, 0.2)",
                    mortgageCompareData[idx]["valid"]
                        ? "rgba(249, 166, 0, 1)"
                        : "rgba(249, 166, 0, 0.2)",
                ][idx],
                backgroundColor: [
                    mortgageCompareData[idx]["valid"]
                        ? "rgba(98, 182, 170, 1)"
                        : "rgba(98, 182, 170, 0.1)",
                    mortgageCompareData[idx]["valid"]
                        ? "rgba(133, 81, 182, 1)"
                        : "rgba(133, 81, 182, 0.1)",
                    mortgageCompareData[idx]["valid"]
                        ? "rgba(249, 166, 0, 1)"
                        : "rgba(249, 166, 0, 0.1)",
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
            (_, i) => `${i}`,
        );
        lineChartInstance.data.datasets = updatedDatasets;
        lineChartInstance.update();
    }
    function updateCashVsMortgageChart() {
        if (!cashVsMortgageChartInstance) return;

        // Use dynamic mortgage duration instead of hardcoded 30
        const years = mortgage_duration + 1;
        let i = 0;
        const updatedDatasets = mortgage_percentage.map(
            (mortgage_percentage, idx) => {
                const data = cashVsMortgageData[idx]["values"];
                wealth_cash_vs_mortgage_installments[i] =
                    cashVsMortgageData[idx]["installment"];
                i += 1;
                return {
                    label:
                        mortgage_percentage == 0
                            ? `Cash`
                            : `${mortgage_percentage}% mutuo`,
                    data,
                    borderColor: [
                        "rgba(178, 178, 178, 1)",
                        "rgba(98, 182, 170, 1)",
                        "rgba(133, 81, 182, 1)",
                        "rgba(249, 166, 0, 1)",
                    ][idx],
                    backgroundColor: [
                        "rgba(178, 178, 178, 0.7)",
                        "rgba(98, 182, 170, 0.7)",
                        "rgba(133, 81, 182, 0.7)",
                        "rgba(249, 166, 0, 0.7)",
                    ][idx],
                    pointRadius: 2, // <<< smaller dots
                    pointHoverRadius: 4, // <<< slightly bigger on hover
                    tension: 0.3,
                };
            },
        );

        wealth_cash_vs_mortgage[0] =
            updatedDatasets[0].data[updatedDatasets[0].data.length - 1];
        wealth_cash_vs_mortgage[1] =
            updatedDatasets[1].data[updatedDatasets[1].data.length - 1];
        wealth_cash_vs_mortgage[2] =
            updatedDatasets[2].data[updatedDatasets[2].data.length - 1];
        wealth_cash_vs_mortgage[3] =
            updatedDatasets[3].data[updatedDatasets[3].data.length - 1];

        cashVsMortgageChartInstance.data.labels = Array.from(
            { length: years },
            (_, i) => `${i}`,
        );
        cashVsMortgageChartInstance.data.datasets = updatedDatasets;
        cashVsMortgageChartInstance.update();
    }
</script>

<div
    class="min-h-screen bg-gradient-to-b from-purple-400 to-[#1e1f25] flex items-start justify-center p-2 sm:p-6 pt-16 sm:pt-20"
>
    <!-- Login/User Button - Top Right -->
    <div class="fixed top-4 right-4 z-50 user-menu-container">
        {#if $isLoading}
            <!-- Loading state -->
            <div
                class="bg-white/20 backdrop-blur-sm px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg"
            >
                <div
                    class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"
                ></div>
            </div>
        {:else if $isAuthenticated && $user}
            <!-- User Menu -->
            <div class="relative">
                <button
                    on:click={toggleUserMenu}
                    class="bg-white text-purple-600 hover:bg-purple-50 px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg transition-all duration-300 font-medium text-sm sm:text-base hover:scale-105 flex items-center gap-2"
                >
                    {#if $user.photoURL}
                        <img
                            src={$user.photoURL}
                            alt="Profile"
                            class="w-6 h-6 rounded-full"
                        />
                    {:else}
                        <div
                            class="w-6 h-6 bg-purple-500 rounded-full flex items-center justify-center text-white text-xs font-bold"
                        >
                            {$user.firstName?.charAt(0) ||
                                $user.email?.charAt(0).toUpperCase()}
                        </div>
                    {/if}
                    <span class="hidden sm:inline">
                        {$user.firstName ||
                            $user.displayName?.split(" ")[0] ||
                            "Utente"}!
                    </span>
                    <span class="sm:hidden">
                        {$user.firstName ||
                            $user.displayName?.split(" ")[0] ||
                            "Utente"}
                    </span>
                    <svg
                        class="w-4 h-4 transition-transform duration-200"
                        class:rotate-180={showUserMenu}
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M6 9l6 6 6-6" />
                    </svg>
                </button>

                <!-- User Dropdown Menu -->
                {#if showUserMenu}
                    <div
                        class="absolute right-0 mt-2 w-64 bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden"
                        transition:slide={{ duration: 200 }}
                    >
                        <!-- User Info -->
                        <div class="p-4 border-b border-gray-100">
                            <div class="flex items-center gap-3">
                                {#if $user.photoURL}
                                    <img
                                        src={$user.photoURL}
                                        alt="Profile"
                                        class="w-10 h-10 rounded-full"
                                    />
                                {:else}
                                    <div
                                        class="w-10 h-10 bg-purple-500 rounded-full flex items-center justify-center text-white font-bold"
                                    >
                                        {$user.firstName?.charAt(0) ||
                                            $user.email
                                                ?.charAt(0)
                                                .toUpperCase()}
                                    </div>
                                {/if}
                                <div>
                                    <p
                                        class="font-semibold text-gray-900 text-sm"
                                    >
                                        {$user.displayName ||
                                            `${$user.firstName || ""} ${$user.lastName || ""}`.trim() ||
                                            "Utente"}
                                    </p>
                                    <p class="text-gray-500 text-xs">
                                        {$user.email}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <!-- Menu Items -->
                        <div class="py-2">
                            <button
                                on:click={() => {
                                    showUserMenu = false; /* Add profile navigation */
                                    push("/getpro");
                                }}
                                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
                            >
                                <svg
                                    class="w-4 h-4"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path
                                        d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"
                                    />
                                    <circle cx="12" cy="7" r="4" />
                                </svg>
                                Get pro!
                            </button>
                            <button
                                on:click={() => {
                                    showUserMenu = false; /* Add settings navigation */
                                }}
                                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
                            >
                                <svg
                                    class="w-4 h-4"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <circle cx="12" cy="12" r="3" />
                                    <path d="M12 1v6m0 6v6m11-7h-6m-6 0H1" />
                                </svg>
                                Impostazioni
                            </button>
                            <hr class="my-1 border-gray-100" />
                            <button
                                on:click={handleLogout}
                                class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center gap-2"
                            >
                                <svg
                                    class="w-4 h-4"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path
                                        d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
                                    />
                                    <polyline points="16,17 21,12 16,7" />
                                    <line x1="21" y1="12" x2="9" y2="12" />
                                </svg>
                                Esci
                            </button>
                        </div>
                    </div>
                {/if}
            </div>
        {:else}
            <!-- Login Button -->
            <button
                on:click={goToLogIn}
                class="bg-white text-purple-600 hover:bg-purple-50 px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg transition-all duration-300 font-medium text-sm sm:text-base hover:scale-105"
            >
                Accedi
            </button>
        {/if}
    </div>
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
    <div
        class="w-full max-w-[1350px] mx-auto"
        in:fly={{ x: -1000, duration: 500, delay: 100 }}
    >
        <!-- Mobile: Stack vertically, Desktop: Side by side -->
        <div
            class="flex flex-col lg:flex-row gap-4 sm:gap-6 min-h-[70vh] items-start"
        >
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
                                    ? () => {
                                          updateData();
                                          updateCashVsMortgage();
                                          updateMortgageCompare();
                                      }
                                    : () => {}}
                                on:change={showCosts &&
                                selectedTab != "mortgage"
                                    ? () => {
                                          updateData();
                                          updateCashVsMortgage();
                                          updateMortgageCompare();
                                      }
                                    : () => {}}
                            />
                            <!-- Separator Line -->
                            <div
                                class="absolute inset-y-0 right-10 w-px bg-white"
                            ></div>
                            <!-- Euro Symbol -->
                            <div
                                class="absolute inset-y-0 right-4 flex items-center text-white font-semibold"
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
                                ? () => {
                                      updateData();
                                      updateCashVsMortgage();
                                      updateMortgageCompare();
                                  }
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
                                ? () => {
                                      updateData();
                                      updateCashVsMortgage();
                                      updateMortgageCompare();
                                  }
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
                                ? () => {
                                      updateData();
                                      updateCashVsMortgage();
                                      updateMortgageCompare();
                                  }
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
                                            ? () => {
                                                  updateData();
                                                  updateCashVsMortgage();
                                                  updateMortgageCompare();
                                              }
                                            : () => {}}
                                        on:mouseup={showCosts
                                            ? () => {
                                                  updateData();
                                                  updateCashVsMortgage();
                                                  updateMortgageCompare();
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
                                                      updateMortgageCompare();
                                                  }
                                                : () => {}}
                                        />
                                        <!-- Separator Line -->
                                        <div
                                            class="absolute inset-y-0 right-10 w-px bg-white"
                                        ></div>
                                        <!-- Euro Symbol -->
                                        <div
                                            class="absolute inset-y-0 right-4 flex items-center text-white font-semibold"
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
                                                    ? () => {
                                                          updateData();
                                                          updateCashVsMortgage();
                                                      }
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
                                                          updateCashVsMortgage();
                                                          updateMortgageCompare();
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
                                <NavBar
                                    bind:selectedTab
                                    pro={$user?.pro || false}
                                    getProFunction={goToPro}
                                />

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
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <div
                                                class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6 mb-6"
                                                transition:slide={{
                                                    duration: 500,
                                                }}
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
                                            <h2
                                                class="font-bold text-white leading-tight text-base sm:text-lg mt-4 sm:mt-2 mb-2 text-center"
                                            >
                                                Cosa vado a <span
                                                    class="font-bold text-purple-400"
                                                    >restituire</span
                                                > ogni anno?
                                            </h2>

                                            <!-- Chart - responsive sizing -->
                                            <div
                                                class="w-full h-[300px] sm:h-[350px] relative"
                                            >
                                                <div
                                                    class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent rounded-lg -z-10"
                                                ></div>

                                                {#if mortgage_amount == 0 || mortgage_amount == null || taeg == 0 || taeg == null || mortgage_duration == 0 || mortgage_duration == null}
                                                    <div
                                                        class="absolute inset-0 flex items-center justify-center text-center px-4 z-10"
                                                    >
                                                        <p
                                                            class="text-gray-500 text-sm sm:text-base"
                                                        >
                                                            Inserisci <strong
                                                                >Importo Mutuo</strong
                                                            >,
                                                            <strong
                                                                >Durata</strong
                                                            >
                                                            e
                                                            <strong>TAEG</strong
                                                            > per visualizzare il
                                                            grafico.
                                                        </p>
                                                    </div>
                                                {/if}
                                                <canvas
                                                    id="interestsBarChart"
                                                    class="absolute top-0 left-0 w-full h-full transition-opacity duration-500 ease-in-out"
                                                    class:opacity-0={mortgage_amount ==
                                                        0 ||
                                                        mortgage_amount ==
                                                            null ||
                                                        taeg == 0 ||
                                                        taeg == null ||
                                                        mortgage_duration ==
                                                            0 ||
                                                        mortgage_duration ==
                                                            null}
                                                    class:opacity-100={mortgage_amount >
                                                        0 &&
                                                        taeg > 0 &&
                                                        mortgage_duration > 0}
                                                ></canvas>
                                            </div>
                                        </div>
                                    </div>
                                {/if}

                                {#if selectedTab == "mortgage_compare"}
                                    <!-- Mortgage comparison section -->
                                    <div transition:fade={{ duration: 500 }}>
                                        <!-- Input section - responsive layout -->
                                        <!-- Simple Mortgage comparison section -->
                                        <div
                                            transition:fade={{ duration: 500 }}
                                        >
                                            <!-- Input section -->
                                            <div
                                                class="flex flex-col items-center gap-4 mb-4"
                                                transition:slide={{
                                                    duration: 500,
                                                }}
                                            >
                                                <!-- Text labels and inputs aligned -->
                                                <div
                                                    class="flex justify-center items-center gap-4"
                                                >
                                                    <!-- First column: Euro input with label -->
                                                    <div
                                                        class="flex flex-col items-center gap-2"
                                                    >
                                                        <h3
                                                            class="font-bold text-white text-sm sm:text-base text-center"
                                                        >
                                                            Quanto guadagno/amo <br
                                                            />
                                                            ogni anno?
                                                        </h3>
                                                        <div class="relative">
                                                            <input
                                                                type="number"
                                                                class="w-32 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-white focus:outline-none"
                                                                min="0"
                                                                step="2500"
                                                                bind:value={
                                                                    yearlySaving
                                                                }
                                                                on:change={showCosts
                                                                    ? updateMortgageCompare
                                                                    : () => {}}
                                                            />
                                                            <span
                                                                class="absolute right-3 top-2 text-white font-semibold"
                                                                >€</span
                                                            >
                                                        </div>
                                                    </div>

                                                    <!-- First Arrow with Icon -->
                                                    <div
                                                        class="flex flex-col items-center gap-1 mt-8"
                                                    >
                                                        <!-- Piggy bank icon for savings -->
                                                        💰
                                                        <!-- Arrow -->
                                                        <svg
                                                            class="w-6 h-6 text-white/70"
                                                            fill="none"
                                                            stroke="currentColor"
                                                            viewBox="0 0 24 24"
                                                        >
                                                            <path
                                                                stroke-linecap="round"
                                                                stroke-linejoin="round"
                                                                stroke-width="2"
                                                                d="M17 8l4 4m0 0l-4 4m4-4H3"
                                                            />
                                                        </svg>
                                                    </div>

                                                    <!-- Second column: Savings rate input -->
                                                    <div
                                                        class="flex flex-col items-center gap-2"
                                                    >
                                                        <h3
                                                            class="font-bold text-white text-sm sm:text-base text-center"
                                                        >
                                                            Di cui riesco/amo a <br
                                                            />risparmiare un...
                                                        </h3>
                                                        <div class="relative">
                                                            <input
                                                                type="number"
                                                                class="w-24 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-white focus:outline-none"
                                                                min="0"
                                                                step="1"
                                                                bind:value={
                                                                    yearlySavingRate
                                                                }
                                                                on:change={showCosts
                                                                    ? updateMortgageCompare
                                                                    : () => {}}
                                                            />
                                                            <span
                                                                class="absolute right-3 top-2 text-white font-semibold"
                                                                >%</span
                                                            >
                                                        </div>
                                                    </div>

                                                    <!-- Second Arrow with Icon -->
                                                    <div
                                                        class="flex flex-col items-center gap-1 mt-8"
                                                    >
                                                        <!-- Growth/chart icon for investment returns -->
                                                        💸
                                                        <!-- Arrow -->
                                                        <svg
                                                            class="w-6 h-6 text-white/70"
                                                            fill="none"
                                                            stroke="currentColor"
                                                            viewBox="0 0 24 24"
                                                        >
                                                            <path
                                                                stroke-linecap="round"
                                                                stroke-linejoin="round"
                                                                stroke-width="2"
                                                                d="M17 8l4 4m0 0l-4 4m4-4H3"
                                                            />
                                                        </svg>
                                                    </div>

                                                    <!-- Third column: Investment growth rate -->
                                                    <div
                                                        class="flex flex-col items-center gap-2"
                                                    >
                                                        <h3
                                                            class="font-bold text-white text-sm sm:text-base text-center"
                                                        >
                                                            Quanto fruttano i
                                                            risparmi <br />ogni
                                                            anno?
                                                        </h3>
                                                        <div class="relative">
                                                            <input
                                                                type="number"
                                                                class="w-24 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-white focus:outline-none"
                                                                min="0"
                                                                step="0.1"
                                                                bind:value={
                                                                    yearlyGrowthgRate
                                                                }
                                                                on:change={showCosts
                                                                    ? updateMortgageCompare
                                                                    : () => {}}
                                                            />
                                                            <span
                                                                class="absolute right-3 top-2 text-white font-semibold"
                                                                >%</span
                                                            >
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>

                                        <!-- Question with reduced top margin -->
                                        <h2
                                            class="font-bold text-white leading-tight text-base sm:text-lg text-center mb-2 mt-3"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            Come cresce il mio/nostro <span
                                                class="font-bold text-purple-400"
                                                >patrimonio</span
                                            >
                                            con
                                            <span
                                                class="font-bold text-purple-400"
                                                >il mutuo selezionato</span
                                            >, ma con diverse durate?
                                        </h2>
                                        <!-- Chart - responsive sizing with subtle enhancement -->
                                        <div
                                            class="w-full h-[300px] sm:h-[300px] lg:h-[350px] mb-6 relative"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <!-- Subtle background enhancement -->
                                            <div
                                                class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent rounded-lg -z-10"
                                            ></div>
                                            {#if mortgage_amount == 0 || mortgage_amount == null || taeg == 0 || taeg == null}
                                                <div
                                                    class="absolute inset-0 flex items-center justify-center text-center px-4 z-10"
                                                >
                                                    <p
                                                        class="text-gray-500 text-sm sm:text-base"
                                                    >
                                                        Inserisci <strong
                                                            >Importo Mutuo</strong
                                                        >
                                                        e
                                                        <strong>TAEG</strong> per
                                                        visualizzare il grafico.
                                                    </p>
                                                </div>
                                            {/if}
                                            <canvas
                                                id="lineChart"
                                                class="absolute top-0 left-0 w-full h-full transition-opacity duration-500 ease-in-out"
                                                class:opacity-0={mortgage_amount ==
                                                    0 ||
                                                    mortgage_amount == null ||
                                                    taeg == 0 ||
                                                    taeg == null}
                                                class:opacity-100={mortgage_amount >
                                                    0 && taeg > 0}
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
                                                        isValid={mortgageCompareData[0][
                                                            "valid"
                                                        ]}
                                                        secondaryNumber={mortgageCompareData[0][
                                                            "installment"
                                                        ] / 12}
                                                        secondaryLabel={"rata"}
                                                        years={30}
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
                                                        class="absolute inset-0 bg-gradient-to-br from-indigo-700/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth[1]}
                                                        name={"Mutuo 20 Anni"}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(133, 81, 182, 1)"}
                                                        isValid={mortgageCompareData[1][
                                                            "valid"
                                                        ]}
                                                        secondaryNumber={mortgageCompareData[1][
                                                            "installment"
                                                        ] / 12}
                                                        secondaryLabel={"rata"}
                                                        years={30}
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
                                                        isValid={mortgageCompareData[2][
                                                            "valid"
                                                        ]}
                                                        secondaryNumber={mortgageCompareData[2][
                                                            "installment"
                                                        ] / 12}
                                                        secondaryLabel={"rata"}
                                                        years={30}
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                {/if}
                                {#if selectedTab == "cash_vs_mortgage"}
                                    <!-- Cash vs mortgage section -->
                                    <div transition:fade={{ duration: 500 }}>
                                        <!-- Input section -->
                                        <div
                                            class="flex flex-col items-center gap-4 mb-4"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <!-- First row: Main financial inputs -->
                                            <div
                                                class="flex justify-center items-center gap-4"
                                            >
                                                <!-- Annual income input -->
                                                <div
                                                    class="flex flex-col items-center gap-2"
                                                >
                                                    <h3
                                                        class="font-bold text-white text-sm sm:text-base text-center"
                                                    >
                                                        Quanto guadagno/amo <br
                                                        />ogni anno?
                                                    </h3>
                                                    <div class="relative">
                                                        <input
                                                            type="number"
                                                            class="w-32 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-white focus:outline-none"
                                                            min="0"
                                                            step="2500"
                                                            bind:value={
                                                                yearlySaving
                                                            }
                                                            on:change={showCosts
                                                                ? updateCashVsMortgage
                                                                : () => {}}
                                                        />
                                                        <span
                                                            class="absolute right-3 top-2 text-white font-semibold"
                                                            >€</span
                                                        >
                                                    </div>
                                                </div>

                                                <!-- Arrow -->
                                                <div
                                                    class="flex flex-col items-center gap-1 mt-8"
                                                >
                                                    💰
                                                    <svg
                                                        class="w-6 h-6 text-white/70"
                                                        fill="none"
                                                        stroke="currentColor"
                                                        viewBox="0 0 24 24"
                                                    >
                                                        <path
                                                            stroke-linecap="round"
                                                            stroke-linejoin="round"
                                                            stroke-width="2"
                                                            d="M17 8l4 4m0 0l-4 4m4-4H3"
                                                        />
                                                    </svg>
                                                </div>

                                                <!-- Savings rate input -->
                                                <div
                                                    class="flex flex-col items-center gap-2"
                                                >
                                                    <h3
                                                        class="font-bold text-white text-sm sm:text-base text-center"
                                                    >
                                                        Di cui riesco/amo a <br
                                                        />risparmiare un...
                                                    </h3>
                                                    <div class="relative">
                                                        <input
                                                            type="number"
                                                            class="w-24 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-white focus:outline-none"
                                                            min="0"
                                                            step="1"
                                                            bind:value={
                                                                yearlySavingRate
                                                            }
                                                            on:change={showCosts
                                                                ? updateCashVsMortgage
                                                                : () => {}}
                                                        />
                                                        <span
                                                            class="absolute right-3 top-2 text-white font-semibold"
                                                            >%</span
                                                        >
                                                    </div>
                                                </div>

                                                <!-- Arrow -->
                                                <div
                                                    class="flex flex-col items-center gap-1 mt-8"
                                                >
                                                    💸
                                                    <svg
                                                        class="w-6 h-6 text-white/70"
                                                        fill="none"
                                                        stroke="currentColor"
                                                        viewBox="0 0 24 24"
                                                    >
                                                        <path
                                                            stroke-linecap="round"
                                                            stroke-linejoin="round"
                                                            stroke-width="2"
                                                            d="M17 8l4 4m0 0l-4 4m4-4H3"
                                                        />
                                                    </svg>
                                                </div>

                                                <!-- Investment growth rate input -->
                                                <div
                                                    class="flex flex-col items-center gap-2"
                                                >
                                                    <h3
                                                        class="font-bold text-white text-sm sm:text-base text-center"
                                                    >
                                                        Quanto fruttano i
                                                        risparmi <br />ogni
                                                        anno?
                                                    </h3>
                                                    <div class="relative">
                                                        <input
                                                            type="number"
                                                            class="w-24 border border-white rounded px-3 py-2 pr-8 text-white bg-transparent focus:border-white focus:outline-none"
                                                            min="0"
                                                            step="0.1"
                                                            bind:value={
                                                                yearlyGrowthgRate
                                                            }
                                                            on:change={showCosts
                                                                ? updateCashVsMortgage
                                                                : () => {}}
                                                        />
                                                        <span
                                                            class="absolute right-3 top-2 text-white font-semibold"
                                                            >%</span
                                                        >
                                                    </div>
                                                </div>
                                            </div>
                                        </div>

                                        <!-- Question -->
                                        <h2
                                            class="font-bold text-white leading-tight text-base sm:text-lg text-center mb-2 mt-3"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            Mi conviene pagare tutto <span
                                                class="font-bold text-purple-400"
                                                >cash</span
                                            >
                                            oppure prendere un
                                            <span
                                                class="font-bold text-purple-400"
                                                >mutuo</span
                                            >?
                                        </h2>

                                        <!-- Chart -->
                                        <div
                                            class="w-full h-[300px] sm:h-[300px] lg:h-[350px] mb-6 relative"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <div
                                                class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent rounded-lg -z-10"
                                            ></div>

                                            {#if taeg == 0 || taeg == null || mortgage_duration == 0 || mortgage_duration == null}
                                                <div
                                                    class="absolute inset-0 flex items-center justify-center text-center px-4 z-10"
                                                >
                                                    <p
                                                        class="text-gray-500 text-sm sm:text-base"
                                                    >
                                                        Inserisci
                                                        <strong
                                                            >Durata mutuo</strong
                                                        >
                                                        e
                                                        <strong>TAEG</strong> per
                                                        visualizzare il grafico.
                                                    </p>
                                                </div>
                                            {/if}
                                            <canvas
                                                id="cashVsMortgageChart"
                                                class="absolute top-0 left-0 w-full h-full transition-opacity duration-500 ease-in-out"
                                                class:opacity-0={taeg == 0 ||
                                                    taeg == null ||
                                                    mortgage_duration == 0 ||
                                                    mortgage_duration == null}
                                                class:opacity-100={taeg > 0 &&
                                                    mortgage_duration > 0}
                                            ></canvas>
                                        </div>

                                        <!-- Summary tiles -->
                                        <div
                                            class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-6"
                                            transition:slide={{ duration: 500 }}
                                        >
                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-gray-400/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth_cash_vs_mortgage[0]}
                                                        name={"Cash"}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(178, 178, 178, 1)"}
                                                        isValid={true}
                                                        years={mortgage_duration}
                                                        secondaryLabel={"rata"}
                                                        secondaryNumber={wealth_cash_vs_mortgage_installments[0]}
                                                    />
                                                </div>
                                            </div>

                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-emerald-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth_cash_vs_mortgage[1]}
                                                        name={`Mutuo ${mortgage_duration} Anni (25%)`}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(98, 182, 170, 1)"}
                                                        isValid={true}
                                                        years={mortgage_duration}
                                                        secondaryLabel={"rata"}
                                                        secondaryNumber={wealth_cash_vs_mortgage_installments[1]}
                                                    />
                                                </div>
                                            </div>

                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-indigo-700/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth_cash_vs_mortgage[2]}
                                                        name={`Mutuo ${mortgage_duration} Anni (50%)`}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(133, 81, 182, 1)"}
                                                        isValid={true}
                                                        years={mortgage_duration}
                                                        secondaryLabel={"rata"}
                                                        secondaryNumber={wealth_cash_vs_mortgage_installments[2]}
                                                    />
                                                </div>
                                            </div>

                                            <div
                                                class="group cursor-pointer transition-transform duration-200 hover:scale-105"
                                            >
                                                <div
                                                    class="relative overflow-hidden rounded-lg"
                                                >
                                                    <div
                                                        class="absolute inset-0 bg-gradient-to-br from-amber-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                                    ></div>
                                                    <ColoredSummaryPrice
                                                        number={wealth_cash_vs_mortgage[3]}
                                                        name={`Mutuo ${mortgage_duration} Anni (80%)`}
                                                        showVal={mortgageInstallment !=
                                                            null}
                                                        delta={null}
                                                        color={"rgba(249, 166, 0, 1)"}
                                                        isValid={true}
                                                        years={mortgage_duration}
                                                        secondaryLabel={"rata"}
                                                        secondaryNumber={wealth_cash_vs_mortgage_installments[3]}
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                {/if}
                                {#if selectedTab == "prices"}
                                    <Prices
                                        bind:showErrorPopup
                                        bind:showRateLimitPopup
                                    />
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
