<script lang="ts">
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    export let showVal: boolean;
    export let number: number;
    export let delta: any;
    export let name: string;
    export let color: string;
    export let isValid: boolean = true;
    export let secondaryNumber: number | null = null;
    export let secondaryLabel: string = "";
    export let years: number | null = null;
</script>

<div class="flex-1 min-w-[160px]">
    <div
        class="rounded-xl p-3 flex-1 text-center transition-all duration-500 border-2"
        class:bg-neutral-900={isValid}
        class:bg-neutral-800={!isValid}
        class:opacity-60={!isValid}
        class:border-dashed={!isValid}
        class:border-solid={isValid}
        class:border-gray-500={!isValid}
        style="border-color: {isValid ? color : '#6b7280'};"
    >
        <h1 class="text-lg font-bold mb-1" class:text-gray-400={!isValid}>
            {name}
        </h1>

        <div class="flex flex-col items-center">
            <NumberFlowGroup
                style="--number-flow-char-height: 0.85em"
                class="flex items-center font-semibold {!isValid
                    ? 'opacity-75'
                    : ''}"
            >
                <NumberFlow
                    value={number ? number : 0}
                    format={{
                        style: "currency",
                        currency: "EUR",
                        maximumFractionDigits: 0,
                    }}
                    locales={"it-IT"}
                    class="font-semibold {isValid
                        ? 'text-[1.5rem]'
                        : 'text-[1.3rem] text-gray-400'}"
                />
                {#if delta != null}
                    <NumberFlow
                        value={showVal ? delta : 0}
                        format={{
                            style: "percent",
                            maximumFractionDigits: 2,
                            minimumFractionDigits: 2,
                            signDisplay: "always",
                        }}
                        class="font-semibold {isValid
                            ? 'text-lg text-red-500'
                            : 'text-sm text-gray-500'}"
                    />
                    <span class="text-xs {!isValid ? 'text-gray-500' : ''}">
                        in più</span
                    >
                {/if}
            </NumberFlowGroup>

            {#if secondaryNumber != null}
                <div class={!isValid ? "opacity-60" : ""}>
                    <div
                        class="text-xs {isValid
                            ? 'text-gray-400'
                            : 'text-gray-600'} mt-1"
                    >
                        {secondaryLabel}
                    </div>
                    <NumberFlow
                        value={secondaryNumber}
                        format={{
                            style: "currency",
                            currency: "EUR",
                            maximumFractionDigits: 0,
                        }}
                        locales={"it-IT"}
                        class="font-medium {isValid
                            ? 'text-base text-gray-300'
                            : 'text-sm text-gray-500'}"
                    />
                </div>
            {/if}

            {#if isValid}
                <div
                    class="flex items-center gap-1 mt-1 min-h-[30px] transition-opacity duration-500 opacity-100"
                >
                    <svg
                        class="w-5 h-5 text-green-500"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <span class="text-xs text-green-500"
                        >Patrimonio fra {years} anni</span
                    >
                </div>
            {/if}
            {#if !isValid}
                <div
                    class="flex items-center gap-1 mt-1 min-h-[30px] transition-opacity duration-500 opacity-100"
                >
                    <svg
                        class="w-5 h-5 text-red-500"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M9.401 3.003c1.155-2 4.043-2 5.197 0l7.355 12.748c1.154 2-.29 4.5-2.599 4.5H4.645c-2.309 0-3.752-2.5-2.598-4.5L9.4 3.003zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <span class="text-xs text-red-400"
                        >Rata del mutuo troppo alta!</span
                    >
                </div>
            {/if}
        </div>
    </div>
</div>
