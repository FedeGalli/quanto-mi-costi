<script lang="ts">
    export let label: any = "";
    export let bind = false;
    export let showTooltip: boolean[];
    export let element: number = 0;

    async function updateShowTooltip() {
        for (let i = 0; i < showTooltip.length; i++) {
            if (i != element) showTooltip[i] = false;
        }
        showTooltip[element] = !showTooltip[element];
    }
</script>

<div class="checkbox-wrapper-42 flex items-center gap-2 w-full">
    <input id={label} type="checkbox" bind:checked={bind} on:change />
    <label class="cbx" for={label}></label>
    <h1 class="text-base font-bold leading-tight">{label}</h1>

    <!-- Push tooltip button to the right -->
    <div class="ml-auto relative group cursor-pointer">
        <button
            class="w-6 h-6 flex items-center justify-center rounded-full hover:text-purple-400 text-white text-xs focus:outline-none transition"
            on:click={async () => updateShowTooltip()}
        >
            <!-- SVG Icon -->
            <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 512 512"
                fill="currentColor"
                class="w-4 h-4"
            >
                <path
                    d="M256 512A256 256 0 1 0 256 0a256 256 0 1 0 0 512zM169.8 165.3c7.9-22.3 29.1-37.3 52.8-37.3l58.3 0c34.9 0 63.1 28.3 63.1 63.1c0 22.6-12.1 43.5-31.7 54.8L280 264.4c-.2 13-10.9 23.6-24 23.6c-13.3 0-24-10.7-24-24l0-13.5c0-8.6 4.6-16.5 12.1-20.8l44.3-25.4c4.7-2.7 7.6-7.7 7.6-13.1c0-8.4-6.8-15.1-15.1-15.1l-58.3 0c-3.4 0-6.4 2.1-7.5 5.3l-.4 1.2c-4.4 12.5-18.2 19-30.6 14.6s-19-18.2-14.6-30.6l.4-1.2zM224 352a32 32 0 1 1 64 0 32 32 0 1 1 -64 0z"
                />
            </svg>
        </button>
    </div>
</div>

<style>
    .checkbox-wrapper-42 input[type="checkbox"] {
        display: none;
        visibility: hidden;
    }

    .checkbox-wrapper-42 {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .checkbox-wrapper-42 .cbx {
        position: relative;
        top: 1px;
        width: 17px;
        height: 17px;
        border: 1px solid #c8ccd4;
        border-radius: 3px;
        vertical-align: middle;
        transition: background 0.1s ease;
        cursor: pointer;
    }
    .checkbox-wrapper-42 .cbx:after {
        content: "";
        position: absolute;
        top: 1px;
        left: 5px;
        width: 5px;
        height: 11px;
        opacity: 0;
        transform: rotate(45deg) scale(0);
        border-right: 2px solid rgb(0, 0, 0);
        border-bottom: 2px solid rgb(0, 0, 0);
        transition: all 0.3s ease;
        transition-delay: 0.15s;
    }
    .checkbox-wrapper-42 input[type="checkbox"]:checked ~ .cbx {
        border-color: transparent;
        background: rgba(186, 126, 254, 1);
        animation: jelly-42 0.6s ease;
    }
    .checkbox-wrapper-42 input[type="checkbox"]:checked ~ .cbx:after {
        opacity: 1;
        transform: rotate(45deg) scale(1);
    }
    @-moz-keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
    @-webkit-keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
    @-o-keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
    @keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
</style>
