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
            Come capire se la vedita è curata da un <span
                class="text-purple-400">agenzia</span
            >?
        </h2>
        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1.</span> Controllare approfonditamente
            l'annuncio dell'immobile
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                Se l’annuncio riporta il nome dell’agenzia (es. Tecnocasa,
                Gabetti, Remax, ecc.) o contatti con email tipo
                “@agenziaimmobiliare.it”, quasi sicuramente è tramite agenzia.
            </li>
            <li>
                Spesso c’è scritto anche “si richiede provvigione” o “trattativa
                tramite agenzia”.
            </li>
        </ul>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">2.</span> Best practice
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                E' sempre buona pratica, quando si contatta l'inserzionista
                dell'annuncio, chiedere esplicitamente se la vendita è fra
                "privati" oppure "tramite agenzia".
            </li>
            <li>
                Chiedere sempre a quanto ammonta la percentuale di provvigione
                dell'agenzia e verificare quando è dovuto il pagamento di tale
                importo.
            </li>
        </ul>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">3.</span> Pagamento dell'importo
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                Attenzione che alcune agenzie vincolano il pagamento della
                provigione durante l'atto preliminare di vendita e non durante
                il rogito. Pertanto se il rogito non dovesse andare a buon fine
                (per esempio la parte venditrice si ritira dalla vendita, oppure
                non viene concesso il mutuo) il pagamento all'agenzia è comunque
                dovuto!
            </li>
        </ul>
    </div>
</div>
