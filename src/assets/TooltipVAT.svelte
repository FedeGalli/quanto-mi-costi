<script lang="ts">
    export let showTooltip: boolean[] = [];
    export let index: number = 0;
</script>

<!-- Tooltip (Overlay) -->
<div
    class="absolute top-0 left-0 w-full h-full bg-[#1e1f25] rounded-2xl shadow-xl p-8 overflow-auto z-20 transition-all transform duration-600 ease-out-[cubic-bezier(0.22, 1, 0.36, 1)]"
    class:opacity-100={showTooltip[index]}
    class:opacity-0={!showTooltip[index]}
    class:-translate-y-500={!showTooltip[index]}
    class:translate-y-0={showTooltip[index]}
>
    <!-- Close Button -->
    <button
        on:click={() => (showTooltip[index] = false)}
        class="absolute top-6 right-6 text-white text-5xl font-bold hover:text-purple-400 transition"
        aria-label="Chiudi tooltip"
    >
        ×
    </button>
    <div class="space-y-6">
        <h2 class="text-4xl font-bold leading-tight">
            Paghi <span class="text-purple-400">l'IVA</span> se...
        </h2>
        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1.</span> Requisiti sull'immobile acquistato
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                Compri l'immobile da una impresa costruttrice o di
                ristrutturazioni entro 5 anni dalla fine dei lavori (costruzione
                o ristrutturazione).
            </li>
            <li>
                Anche dopo i 5 anni, ma solo se l'impresa sceglie di
                assoggettare la vendita a IVA (è facoltativo, ma deve essere
                indicato nell’atto).
            </li>
            <li>
                Il pagamento dell'IVA, a patto di rispttare le condizione appena
                citate, è comunque dovuto anche in caso di utilizzo di un
                intemediario come per esempio una agenzia immobiliare.
            </li>
        </ul>
    </div>
</div>
