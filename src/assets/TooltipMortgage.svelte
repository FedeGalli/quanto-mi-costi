<script lang="ts">
    export let showTooltip: boolean = false;
</script>

<!-- Tooltip (Overlay) -->
<div
    class="absolute top-0 left-0 w-full h-full bg-[#1e1f25] rounded-2xl shadow-xl p-8 overflow-auto z-20 transition-all transform duration-600 ease-out-[cubic-bezier(0.22, 1, 0.36, 1)]"
    class:opacity-100={showTooltip}
    class:opacity-0={!showTooltip}
    class:-translate-y-500={!showTooltip}
    class:translate-y-0={showTooltip}
>
    <div class="space-y-6">
        <h2 class="text-4xl font-bold leading-tight">
            Informazioni principali per un <span class="text-purple-400"
                >mutuo</span
            >...
        </h2>
        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1.</span> Documentazione richiesta
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                Per lavoratori dipendenti:
                <ul
                    class="list-disc pl-6 mt-2 text-base font-bold leading-tight space-y-2"
                >
                    <li>
                        Certificazione unica (CU), che rappresenza i redditi
                        percepiti per l'anno precendente.
                    </li>
                    <li>Dalle 3 alle 6 ultime buste paga.</li>
                    <li>
                        Estratto conto corrente e lista movimenti degli ultimi 6
                        mesi.
                    </li>
                </ul>
            </li>
            <li>
                Per lavoratori autonomi (partite IVA):
                <ul
                    class="list-disc pl-6 mt-2 text-base font-bold leading-tight space-y-2"
                >
                    <li>Ultimi due modelli unici per persone fisiche.</li>
                    <li>
                        Registrazione telematica degli stessi modelli
                        all’Agenzia delle Entrate.
                    </li>
                    <li>
                        Estratto conto corrente e lista movimenti degli ultimi 6
                        mesi.
                    </li>
                    <li>
                        Visura camerale, che contiene tutti i dati delle imprese
                        iscritte al registro delle imprese.
                    </li>
                </ul>
            </li>
            <li>
                Documenti sull'immobile:
                <ul
                    class="list-disc pl-6 mt-2 text-base font-bold leading-tight space-y-2"
                >
                    <li>Atto preliminare di acquisto.</li>
                    <li>Visura catastale storica.</li>
                    <li>Planimetria catastale.</li>
                    <li>Attestato di prestazione energetica (APE).</li>
                    <li>Atto di provenienza.</li>
                </ul>
            </li>
        </ul>
    </div>
</div>
