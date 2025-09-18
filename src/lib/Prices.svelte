<script lang="ts">
    import { fade, slide } from "svelte/transition";
    import { onMount, tick } from "svelte";
    import * as Chart from "chart.js";
    import CustomButton from "../assets/CustomButton.svelte";

    let selectedComune: string = "";
    let selectedZona: string = "";
    let selectedTipologia: string = "";
    let selectedStato: string = "";
    let selectedMq: number = 0;

    let comuniOptions: any[] = [];
    let comuniInfo: any[] = [];
    let zoneOptions: any[] = [];
    let tipologieOptions: any[] = [];
    let statoOptions: any[] = [];
    let priceVolumesData: any = undefined;

    let chartType = "prices";
    let volumeChart: any;
    let priceChart: any;
    let volumeChartInstance: any;
    let priceChartInstance: any;

    export let apiURL: string = "";

    let validationErrors = {
        comune: false,
        zona: false,
        tipologia: false,
        stato: false,
        mq: false,
    };

    let validationError = "";
    export let showRateLimitPopup = false;
    export let showErrorPopup = false;

    $: if (priceVolumesData && chartType) {
        // Use tick() to wait for DOM updates
        tick().then(() => {
            createChart(chartType);
        });
    }

    function saveStateToLocalStorage() {
        const state = {
            selectedComune,
            selectedZona,
            selectedTipologia,
            selectedStato,
            selectedMq,
            comuniOptions,
            comuniInfo,
            zoneOptions,
            tipologieOptions,
            statoOptions,
            priceVolumesData,
            chartType,
        };
        localStorage.setItem("appState", JSON.stringify(state));
    }
    function loadStateFromLocalStorage() {
        let savedState = localStorage.getItem("appState");
        if (savedState) {
            try {
                let state = JSON.parse(savedState);

                // Restore all variables
                selectedComune = state.selectedComune;
                selectedZona = state.selectedZona;
                selectedTipologia = state.selectedTipologia;
                selectedStato = state.selectedStato;
                selectedMq = state.selectedMq;
                comuniOptions = state.comuniOptions;
                comuniInfo = state.comuniInfo;
                zoneOptions = state.zoneOptions;
                tipologieOptions = state.tipologieOptions;
                statoOptions = state.statoOptions;
                priceVolumesData = state.priceVolumesData;
                chartType = state.chartType;
            } catch (error) {
                console.error("Error loading state from localStorage:", error);
            }
        }
    }

    function createChart(type: string) {
        if (!priceVolumesData) return;

        // Destroy existing chart instances first
        if (volumeChartInstance) {
            volumeChartInstance.destroy();
            volumeChartInstance = null;
        }
        if (priceChartInstance) {
            priceChartInstance.destroy();
            priceChartInstance = null;
        }

        // Create the appropriate chart based on type
        if (type === "volume" && volumeChart) {
            const volumeData =
                priceVolumesData.volume_trend.datasets[0].data.map((v: any) =>
                    parseFloat(v.replace(",", ".")),
                );

            volumeChartInstance = new Chart.Chart(volumeChart, {
                type: "line",
                data: {
                    labels: priceVolumesData.volume_trend.labels,
                    datasets: [
                        {
                            label: "Volume Compravendite",
                            data: volumeData,
                            borderColor: "#a855f7",
                            backgroundColor: "rgba(168, 85, 247, 0.1)",
                            borderWidth: 3,
                            fill: true,
                            tension: 0.4,
                            pointBackgroundColor: "#a855f7",
                            pointBorderColor: "#ffffff",
                            pointBorderWidth: 1,
                            pointRadius: 3,
                            pointHoverRadius: 8,
                        },
                    ],
                },
                options: {
                    responsive: true,
                    interaction: {
                        intersect: false,
                        mode: "index",
                    },
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            labels: {
                                color: "#ffffff",
                                font: { size: 14 },
                            },
                        },
                    },
                    scales: {
                        x: {
                            ticks: { color: "#ffffff" },
                            grid: { color: "rgba(255, 255, 255, 0.1)" },
                        },
                        y: {
                            ticks: { color: "#ffffff" },
                            grid: { color: "rgba(255, 255, 255, 0.1)" },
                        },
                    },
                },
            });
        } else if (type === "prices" && priceChart) {
            priceChartInstance = new Chart.Chart(priceChart, {
                type: "line",
                data: {
                    labels: priceVolumesData.price_trend.labels,
                    datasets: [
                        {
                            label: "Prezzo Minimo m² (€)",
                            data: priceVolumesData.price_trend.datasets[0].data,
                            borderColor: "#10b981",
                            backgroundColor: "rgba(16, 185, 129, 0.1)",
                            borderWidth: 3,
                            fill: false,
                            tension: 0.4,
                            pointBackgroundColor: "#10b981",
                            pointBorderColor: "#ffffff",
                            pointBorderWidth: 1,
                            pointRadius: 3,
                            pointHoverRadius: 8,
                        },
                        {
                            label: "Prezzo Massimo m² (€)",
                            data: priceVolumesData.price_trend.datasets[1].data,
                            borderColor: "#ef4444",
                            backgroundColor: "rgba(239, 68, 68, 0.1)",
                            borderWidth: 3,
                            fill: false,
                            tension: 0.4,
                            pointBackgroundColor: "#ef4444",
                            pointBorderColor: "#ffffff",
                            pointBorderWidth: 1,
                            pointRadius: 3,
                            pointHoverRadius: 8,
                        },
                    ],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    interaction: {
                        intersect: false,
                        mode: "index",
                    },
                    plugins: {
                        legend: {
                            labels: {
                                color: "#ffffff",
                                font: { size: 14 },
                            },
                        },
                    },
                    scales: {
                        x: {
                            ticks: { color: "#ffffff" },
                            grid: { color: "rgba(255, 255, 255, 0.1)" },
                        },
                        y: {
                            ticks: { color: "#ffffff" },
                            grid: { color: "rgba(255, 255, 255, 0.1)" },
                            grace: "10%",
                        },
                    },
                },
            });
        }
    }

    // Dropdown visibility state
    let showComuniDropdown: boolean = false;

    function getMunicipalitiesString(): string {
        return apiURL + "/get-municipalities-list";
    }
    function getMunicipalitiesInfoString(com: string): string {
        return apiURL + "/get-municipalities-info?com=" + com;
    }
    function getPriceVolumeDataString(
        com: string,
        zona: string,
        tipo: string,
        stato: string,
        mq: number,
    ): string {
        return (
            apiURL +
            "/get-price-volume-data?com=" +
            com +
            "&zone=" +
            zona +
            "&type=" +
            tipo +
            "&state=" +
            stato +
            "&mq=" +
            mq
        );
    }

    async function getMunicipalitiesData() {
        if (comuniOptions.length != 0) return;
        const apiStringUrl: string = getMunicipalitiesString();
        let responseCode: string = "200";

        try {
            const response = await fetch(apiStringUrl);

            if (!response.ok || response.status == 429) {
                throw new Error(`${response.status}`);
            }
            const data = await response.json();
            comuniOptions = data["DATA"];
            saveStateToLocalStorage();
        } catch (error: any) {
            responseCode = error.message;
            if (responseCode == "429") {
                showRateLimitPopup = true;
            } else {
                showErrorPopup = true;
            }
        }
    }

    async function getPriceVolumeData(
        com: string,
        zona: string,
        tipo: string,
        stato: string,
        mq: number,
    ) {
        const apiStringUrl: string = getPriceVolumeDataString(
            com,
            zona,
            tipo,
            stato,
            mq,
        );
        let responseCode: string = "200";

        try {
            const response = await fetch(apiStringUrl);

            if (!response.ok || response.status == 429) {
                throw new Error(`${response.status}`);
            }
            const data = await response.json();
            priceVolumesData = data;
            saveStateToLocalStorage();
        } catch (error: any) {
            responseCode = error.message;
            if (responseCode == "429") {
                showRateLimitPopup = true;
            } else {
                showErrorPopup = true;
            }
        }
    }

    async function getMunicipalitiesInfoData(com: string) {
        const apiStringUrl: string = getMunicipalitiesInfoString(com);
        let responseCode: string = "200";

        try {
            const response = await fetch(apiStringUrl);

            if (!response.ok || response.status == 429) {
                throw new Error(`${response.status}`);
            }
            const data = await response.json();
            comuniInfo = data["DATA"];
            saveStateToLocalStorage();
        } catch (error: any) {
            responseCode = error.message;
            if (responseCode == "429") {
                showRateLimitPopup = true;
            } else {
                showErrorPopup = true;
            }
        }
    }

    function getZoneOptions() {
        zoneOptions = [];
        const zoneSet = new Set<string>(); // or Set<any> if not string
        comuniInfo.forEach((prop) => {
            zoneSet.add(prop[1]);
        });
        zoneOptions = Array.from(zoneSet);
    }

    function getTypeOptions(zona: string) {
        tipologieOptions = [];
        const tipologieSet = new Set<string>(); // or Set<any> if not string
        comuniInfo.forEach((prop) => {
            if (prop[1] == zona) {
                tipologieSet.add(prop[2]);
            }
        });
        tipologieOptions = Array.from(tipologieSet);
    }

    function getStatoOptions(zona: string) {
        statoOptions = [];
        const statoSet = new Set<string>(); // or Set<any> if not string
        comuniInfo.forEach((prop) => {
            if (prop[2] == zona) {
                statoSet.add(prop[3]);
            }
        });
        statoOptions = Array.from(statoSet);
    }

    let filteredComuni: string[] = [];

    function filterComuni(): void {
        if (selectedComune.length === 0) {
            filteredComuni = comuniOptions.slice(0, 10); // Show first 10 by default
        } else {
            const searchTerm = selectedComune.toLowerCase();

            // Filter and sort results by relevance
            filteredComuni = comuniOptions
                .filter((comune: string) =>
                    comune.toLowerCase().includes(searchTerm),
                )
                .sort((a: string, b: string) => {
                    const aLower = a.toLowerCase();
                    const bLower = b.toLowerCase();

                    // Priority 1: Exact match (case insensitive)
                    if (aLower === searchTerm) return -1;
                    if (bLower === searchTerm) return 1;

                    // Priority 2: Starts with search term
                    const aStarts = aLower.startsWith(searchTerm);
                    const bStarts = bLower.startsWith(searchTerm);
                    if (aStarts && !bStarts) return -1;
                    if (bStarts && !aStarts) return 1;

                    // Priority 3: Alphabetical order for remaining matches
                    return a.localeCompare(b);
                })
                .slice(0, 10); // Limit to 10 results for performance
        }
    }

    function comuneChange() {
        selectedZona = "";
        selectedTipologia = "";
        selectedStato = "";
        zoneOptions = [];
        tipologieOptions = [];
        statoOptions = [];
    }

    // Selection functions
    async function selectComune(comune: string) {
        selectedComune = comune;
        showComuniDropdown = false;
        comuneChange();
        await getMunicipalitiesInfoData(comune);
        getZoneOptions();
    }

    function zonaChange() {
        selectedTipologia = "";
        selectedStato = "";
    }

    function selectZona(zona: string): void {
        zonaChange();
        getTypeOptions(zona);
    }

    function tipologiaChange() {
        selectedStato = "";
    }

    function selectTipologia(zona: string): void {
        tipologiaChange();
        getStatoOptions(zona);
    }

    onMount(async () => {
        loadStateFromLocalStorage();
        await getMunicipalitiesData();
        filterComuni();
        // Register Chart.js components
        Chart.Chart.register(
            Chart.LineElement,
            Chart.PointElement,
            Chart.CategoryScale,
            Chart.LinearScale,
            Chart.Title,
            Chart.Tooltip,
            Chart.Legend,
            Chart.Filler,
        );
    });
</script>

<!-- Prices section -->
<div transition:fade={{ duration: 500 }}>
    <div
        class="flex flex-col items-center gap-4 mb-4"
        transition:slide={{ duration: 500 }}
    >
        {#if !priceVolumesData}
            <!-- Input section - Only show when no data -->
            <!-- Question -->
            <h2
                class="font-bold text-white leading-tight text-base sm:text-lg text-center mb-4"
            >
                Quanto <span class="font-bold text-purple-400">costano</span>
                le case nella mia
                <span class="font-bold text-purple-400">zona</span>?
            </h2>

            <!-- Form Fields - Stacked Layout for narrow panel -->
            <div class="w-full space-y-4 max-w-sm">
                <!-- Comune (Autocomplete) - Step 1 -->
                <div class="flex flex-col gap-2">
                    <div class="flex items-center gap-2">
                        <span
                            class="text-purple-400 font-bold text-xs bg-purple-400/20 px-2 py-1 rounded"
                            >1</span
                        >
                        <h3 class="font-bold text-white text-sm text-left">
                            Comune
                        </h3>
                        {#if selectedComune}
                            <svg
                                class="w-4 h-4 text-green-400"
                                fill="currentColor"
                                viewBox="0 0 20 20"
                            >
                                <path
                                    fill-rule="evenodd"
                                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                    clip-rule="evenodd"
                                />
                            </svg>
                        {/if}
                    </div>
                    <div class="relative">
                        <input
                            type="text"
                            class="w-full border rounded-lg px-4 py-3 text-white bg-transparent focus:outline-none placeholder-gray-300 text-sm {validationErrors.comune
                                ? 'border-red-500 focus:border-red-500 focus:ring-2 focus:ring-red-500/20'
                                : 'border-white focus:border-purple-400 focus:ring-2 focus:ring-purple-400/20'}"
                            placeholder="Cerca comune..."
                            bind:value={selectedComune}
                            on:input={filterComuni}
                            on:focus={() => (showComuniDropdown = true)}
                            on:blur={() =>
                                setTimeout(
                                    () => (showComuniDropdown = false),
                                    200,
                                )}
                        />
                        {#if showComuniDropdown && filteredComuni.length > 0}
                            <div
                                class="absolute top-full left-0 right-0 bg-white border border-gray-300 rounded-b-lg max-h-40 overflow-y-auto z-20 shadow-lg mt-1"
                            >
                                {#each filteredComuni as comune}
                                    <div
                                        class="px-4 py-3 hover:bg-gray-100 cursor-pointer text-black text-sm border-b border-gray-100 last:border-b-0"
                                        on:click={() => selectComune(comune)}
                                    >
                                        {comune}
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </div>

                <!-- Two-column row for Zona and Tipologia -->
                <div class="grid grid-cols-2 gap-3">
                    <!-- Zona (Dropdown) - Step 2 -->
                    <div class="flex flex-col gap-2">
                        <div class="flex items-center gap-2">
                            <span
                                class="text-purple-400 font-bold text-xs bg-purple-400/20 px-2 py-1 rounded"
                                >2</span
                            >
                            <h3
                                class="font-bold text-sm text-left {selectedComune
                                    ? 'text-white'
                                    : 'text-gray-400'}"
                            >
                                Zona
                            </h3>
                            {#if selectedZona}
                                <svg
                                    class="w-4 h-4 text-green-400"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                >
                                    <path
                                        fill-rule="evenodd"
                                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                        clip-rule="evenodd"
                                    />
                                </svg>
                            {/if}
                        </div>
                        <div class="relative">
                            <select
                                class="w-full border rounded-lg px-3 py-3 bg-transparent focus:outline-none appearance-none cursor-pointer text-sm pr-8 {selectedComune
                                    ? validationErrors.zona
                                        ? 'border-red-500 text-white focus:border-red-500 focus:ring-2 focus:ring-red-500/20'
                                        : 'border-white text-white focus:border-purple-400 focus:ring-2 focus:ring-purple-400/20'
                                    : 'border-gray-500 text-gray-400 cursor-not-allowed'}"
                                bind:value={selectedZona}
                                disabled={!selectedComune}
                                on:change={() => {
                                    if (selectedComune) {
                                        selectZona(selectedZona);
                                    }
                                }}
                            >
                                <option value="" class="text-black bg-white">
                                    {selectedComune
                                        ? "Seleziona"
                                        : "Prima scegli comune"}
                                </option>
                                {#each zoneOptions as zona}
                                    <option
                                        value={zona}
                                        class="text-black bg-white"
                                        >{zona}</option
                                    >
                                {/each}
                            </select>
                            <div
                                class="absolute right-2 top-1/2 transform -translate-y-1/2 pointer-events-none"
                            >
                                <svg
                                    class="w-4 h-4 {selectedComune
                                        ? 'text-white'
                                        : 'text-gray-400'}"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M19 9l-7 7-7-7"
                                    ></path>
                                </svg>
                            </div>
                        </div>
                    </div>

                    <!-- Tipologia (Dropdown) - Step 3 -->
                    <div class="flex flex-col gap-2">
                        <div class="flex items-center gap-2">
                            <span
                                class="text-purple-400 font-bold text-xs bg-purple-400/20 px-2 py-1 rounded"
                                >3</span
                            >
                            <h3
                                class="font-bold text-sm text-left {selectedZona
                                    ? 'text-white'
                                    : 'text-gray-400'}"
                            >
                                Tipologia
                            </h3>
                            {#if selectedTipologia}
                                <svg
                                    class="w-4 h-4 text-green-400"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                >
                                    <path
                                        fill-rule="evenodd"
                                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                        clip-rule="evenodd"
                                    />
                                </svg>
                            {/if}
                        </div>
                        <div class="relative">
                            <select
                                class="w-full border rounded-lg px-3 py-3 bg-transparent focus:outline-none appearance-none cursor-pointer text-sm pr-8 {selectedZona
                                    ? validationErrors.tipologia
                                        ? 'border-red-500 text-white focus:border-red-500 focus:ring-2 focus:ring-red-500/20'
                                        : 'border-white text-white focus:border-purple-400 focus:ring-2 focus:ring-purple-400/20'
                                    : 'border-gray-500 text-gray-400 cursor-not-allowed'}"
                                bind:value={selectedTipologia}
                                disabled={!selectedZona}
                                on:change={() => {
                                    if (selectedZona) {
                                        selectTipologia(selectedTipologia);
                                    }
                                }}
                            >
                                <option value="" class="text-black bg-white">
                                    {selectedZona
                                        ? "Seleziona"
                                        : "Prima scegli zona"}
                                </option>
                                {#each tipologieOptions as tipologia}
                                    <option
                                        value={tipologia}
                                        class="text-black bg-white"
                                        >{tipologia}</option
                                    >
                                {/each}
                            </select>
                            <div
                                class="absolute right-2 top-1/2 transform -translate-y-1/2 pointer-events-none"
                            >
                                <svg
                                    class="w-4 h-4 {selectedZona
                                        ? 'text-white'
                                        : 'text-gray-400'}"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M19 9l-7 7-7-7"
                                    ></path>
                                </svg>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Two-column row for Stato and m² -->
                <div class="grid grid-cols-2 gap-3">
                    <!-- Stato (Dropdown) - Step 4 -->
                    <div class="flex flex-col gap-2">
                        <div class="flex items-center gap-2">
                            <span
                                class="text-purple-400 font-bold text-xs bg-purple-400/20 px-2 py-1 rounded"
                                >4</span
                            >
                            <h3
                                class="font-bold text-sm text-left {selectedTipologia
                                    ? 'text-white'
                                    : 'text-gray-400'}"
                            >
                                Stato
                            </h3>
                            {#if selectedStato}
                                <svg
                                    class="w-4 h-4 text-green-400"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                >
                                    <path
                                        fill-rule="evenodd"
                                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                        clip-rule="evenodd"
                                    />
                                </svg>
                            {/if}
                        </div>
                        <div class="relative">
                            <select
                                class="w-full border rounded-lg px-3 py-3 bg-transparent focus:outline-none appearance-none cursor-pointer text-sm pr-8 {selectedTipologia
                                    ? validationErrors.stato
                                        ? 'border-red-500 text-white focus:border-red-500 focus:ring-2 focus:ring-red-500/20'
                                        : 'border-white text-white focus:border-purple-400 focus:ring-2 focus:ring-purple-400/20'
                                    : 'border-gray-500 text-gray-400 cursor-not-allowed'}"
                                bind:value={selectedStato}
                                disabled={!selectedTipologia}
                            >
                                <option value="" class="text-black bg-white">
                                    {selectedTipologia
                                        ? "Seleziona"
                                        : "Prima scegli tipologia"}
                                </option>
                                {#each statoOptions as stato}
                                    <option
                                        value={stato}
                                        class="text-black bg-white"
                                        >{stato}</option
                                    >
                                {/each}
                            </select>
                            <div
                                class="absolute right-2 top-1/2 transform -translate-y-1/2 pointer-events-none"
                            >
                                <svg
                                    class="w-4 h-4 {selectedTipologia
                                        ? 'text-white'
                                        : 'text-gray-400'}"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M19 9l-7 7-7-7"
                                    ></path>
                                </svg>
                            </div>
                        </div>
                    </div>

                    <!-- m² (Number input) - Step 5 -->
                    <div class="flex flex-col gap-2">
                        <div class="flex items-center gap-2">
                            <span
                                class="text-purple-400 font-bold text-xs bg-purple-400/20 px-2 py-1 rounded"
                                >5</span
                            >
                            <h3
                                class="font-bold text-sm text-left {selectedStato
                                    ? 'text-white'
                                    : 'text-gray-400'}"
                            >
                                m²
                            </h3>
                            {#if selectedMq}
                                <svg
                                    class="w-4 h-4 text-green-400"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                >
                                    <path
                                        fill-rule="evenodd"
                                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                        clip-rule="evenodd"
                                    />
                                </svg>
                            {/if}
                        </div>
                        <input
                            type="number"
                            class="w-full border rounded-lg px-4 py-3 bg-transparent focus:outline-none text-sm {selectedStato
                                ? validationErrors.mq
                                    ? 'border-red-500 text-white focus:border-red-500 focus:ring-2 focus:ring-red-500/20'
                                    : 'border-white text-white focus:border-purple-400 focus:ring-2 focus:ring-purple-400/20'
                                : 'border-gray-500 text-gray-400 cursor-not-allowed'}"
                            min="0"
                            max="500"
                            placeholder={selectedStato
                                ? "120"
                                : "Prima scegli stato"}
                            bind:value={selectedMq}
                            disabled={!selectedStato}
                            on:input={() => {
                                if (selectedMq < 0) selectedMq = 0;
                            }}
                        />
                    </div>
                </div>

                <!-- Validation Error Message -->
                {#if validationError}
                    <div
                        class="bg-red-500/20 border border-red-500 rounded-lg p-3 text-red-200 text-sm text-center"
                        transition:slide={{ duration: 500 }}
                    >
                        <div class="flex items-center justify-center gap-2">
                            <svg
                                class="w-4 h-4 text-red-400"
                                fill="currentColor"
                                viewBox="0 0 20 20"
                            >
                                <path
                                    fill-rule="evenodd"
                                    d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                                    clip-rule="evenodd"
                                />
                            </svg>
                            <span>{validationError}</span>
                        </div>
                    </div>
                {/if}

                <!-- Calculate Button -->
                <div class="pt-2">
                    <CustomButton
                        title={"Calcola"}
                        onClick={() => {
                            // Reset previous validation errors
                            validationErrors = {
                                comune: false,
                                zona: false,
                                tipologia: false,
                                stato: false,
                                mq: false,
                            };
                            validationError = "";

                            // Check each field and mark errors
                            let hasErrors = false;
                            let missingFields = [];

                            if (!selectedComune) {
                                validationErrors.comune = true;
                                missingFields.push("Comune");
                                hasErrors = true;
                            }
                            if (!selectedZona) {
                                validationErrors.zona = true;
                                missingFields.push("Zona");
                                hasErrors = true;
                            }
                            if (!selectedTipologia) {
                                validationErrors.tipologia = true;
                                missingFields.push("Tipologia");
                                hasErrors = true;
                            }
                            if (!selectedStato) {
                                validationErrors.stato = true;
                                missingFields.push("Stato");
                                hasErrors = true;
                            }
                            if (!selectedMq) {
                                validationErrors.mq = true;
                                missingFields.push("m²");
                                hasErrors = true;
                            }

                            // Set error message based on missing fields
                            if (hasErrors) {
                                if (missingFields.length === 1) {
                                    validationError = `Completa il campo: ${missingFields[0]}`;
                                } else {
                                    validationError = `Completa i campi: ${missingFields.join(", ")}`;
                                }

                                // Auto-clear error message after 4 seconds
                                setTimeout(() => {
                                    validationError = "";
                                    validationErrors = {
                                        comune: false,
                                        zona: false,
                                        tipologia: false,
                                        stato: false,
                                        mq: false,
                                    };
                                }, 4000);
                            }

                            // Trigger reactivity
                            validationErrors = validationErrors;

                            // If no errors, proceed with calculation
                            if (!hasErrors) {
                                getPriceVolumeData(
                                    selectedComune,
                                    selectedZona,
                                    selectedTipologia,
                                    selectedStato,
                                    selectedMq,
                                );
                            }
                        }}
                        class="w-full"
                    >
                        CALCOLA
                    </CustomButton>
                </div>
            </div>
        {:else}
            <!-- Data Display Section - Takes over main view -->

            <div class="w-full" transition:slide={{ duration: 500 }}>
                <!-- Header with back button -->
                <div class="flex items-center justify-between mb-6">
                    <h2 class="font-bold text-white text-xl">
                        Risultati Analisi
                    </h2>
                    <button
                        class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg text-sm transition-colors"
                        on:click={() => {
                            // Reset to show form again
                            priceVolumesData = null;
                        }}
                    >
                        Nuova Ricerca
                    </button>
                </div>

                <!-- Current Stats Cards -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
                    <!-- Price Range Card -->
                    <div
                        class="bg-neutral-900 border border-purple-400 rounded-xl p-5 flex-1 text-center"
                    >
                        <div class="space-y-4">
                            <!-- Price Range -->
                            <div class="text-center mb-6">
                                <h3
                                    class="text-purple-300 text-sm font-medium mb-2"
                                >
                                    Prezzo Medio Immobile
                                </h3>
                                <div class="text-white text-5xl font-bold mb-1">
                                    {(
                                        ((priceVolumesData.current_max_price +
                                            priceVolumesData.current_min_price) /
                                            2) *
                                        selectedMq
                                    ).toLocaleString("it-IT")} €
                                </div>
                                <div class="text-gray-400 text-sm">
                                    {selectedMq} m² • €{(
                                        (priceVolumesData.current_max_price +
                                            priceVolumesData.current_min_price) /
                                        2
                                    ).toLocaleString("it-IT")}/m²
                                </div>
                            </div>

                            <!-- Visual Range -->
                            <div class="space-y-3">
                                <div class="flex justify-between text-sm">
                                    <span class="text-green-400 font-semibold"
                                        >{(
                                            priceVolumesData.current_min_price *
                                            selectedMq
                                        ).toLocaleString("it-IT")} €</span
                                    >
                                    <span class="text-red-400 font-semibold"
                                        >{(
                                            priceVolumesData.current_max_price *
                                            selectedMq
                                        ).toLocaleString("it-IT")} €</span
                                    >
                                </div>
                                <div
                                    class="w-full bg-gray-700 rounded-full h-3 relative"
                                >
                                    <div
                                        class="bg-gradient-to-r from-green-400 via-purple-400 to-red-400 h-3 rounded-full"
                                    ></div>
                                    <!-- Simple arrow pointing down -->
                                    <div
                                        class="absolute top-0 left-1/2 transform -translate-x-1/2 -translate-y-full"
                                    >
                                        <div class="text-white text-3xl">↓</div>
                                    </div>
                                </div>
                                <div class="text-center text-xs text-gray-400">
                                    Range di prezzo
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Volume Card -->
                    <div
                        class="bg-neutral-900 border border-purple-400 rounded-xl p-5 flex-1 text-center"
                    >
                        <!-- Market Size Header -->
                        <div class="text-center mb-6">
                            <h3
                                class="text-purple-300 text-sm font-medium mb-2"
                            >
                                Dimensione del Mercato<br />{selectedComune}
                            </h3>

                            <div
                                class="text-5xl font-bold mb-1"
                                class:text-red-500={priceVolumesData.market_size ===
                                    "S"}
                                class:text-orange-400={priceVolumesData.market_size ===
                                    "M"}
                                class:text-yellow-400={priceVolumesData.market_size ===
                                    "L"}
                                class:text-green-400={priceVolumesData.market_size ===
                                    "XL"}
                                class:text-blue-400={priceVolumesData.market_size ===
                                    "XXL"}
                            >
                                {priceVolumesData.market_size}
                            </div>
                        </div>

                        <!-- Volume Details -->
                        <div class="text-gray-400 text-sm mt-1">
                            Volume <br />{priceVolumesData.mq_range}
                        </div>
                        <div class="text-white text-2xl font-bold mt-2">
                            {priceVolumesData.current_volume_mq.toLocaleString(
                                "it-IT",
                            )}
                        </div>
                    </div>
                </div>

                <h2 class="font-bold text-white text-xl">Andamento storico</h2>
                <!-- Chart Section -->
                <div
                    class="bg-gradient-to-br from-white/5 to-transparent rounded-lg"
                >
                    <!-- Chart Toggle Buttons -->
                    <div class="flex justify-start mb-6">
                        <div class="bg-[#1e1f25] p-1 rounded-br-lg">
                            <button
                                class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
                                class:bg-purple-400={chartType === "prices"}
                                class:text-white={chartType === "prices"}
                                class:text-gray-300={chartType !== "prices"}
                                on:click={() => (chartType = "prices")}
                            >
                                Prezzi
                            </button>
                            <button
                                class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
                                class:bg-purple-400={chartType === "volume"}
                                class:text-white={chartType === "volume"}
                                class:text-gray-300={chartType !== "volume"}
                                on:click={() => (chartType = "volume")}
                            >
                                Volume
                            </button>
                        </div>
                    </div>

                    <!-- Chart Container -->
                    <div class="h-80 relative">
                        {#if chartType === "volume"}
                            <canvas
                                bind:this={volumeChart}
                                class="w-full h-full"
                            ></canvas>
                        {:else if chartType === "prices"}
                            <canvas bind:this={priceChart} class="w-full h-full"
                            ></canvas>
                        {/if}
                    </div>
                </div>
            </div>
        {/if}
    </div>
</div>
