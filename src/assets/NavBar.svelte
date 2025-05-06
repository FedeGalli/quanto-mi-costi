<script>
    export let selectedTab = "summary";
</script>

<div class="tabs">
    <button
        class="tab {selectedTab === 'summary' ? 'active' : ''}"
        on:click={() => (selectedTab = "summary")}
        type="button"
    >
        📊
        <span><h1 class="text-xs font-bold leading-tight">Summary</h1></span>
    </button>
    <button
        class="tab {selectedTab === 'base' ? 'active' : ''}"
        on:click={() => (selectedTab = "base")}
        type="button"
    >
        📦
        <span><h1 class="text-xs font-bold leading-tight">Base Costs</h1></span>
    </button>

    <button
        class="tab {selectedTab === 'mortgage' ? 'active' : ''}"
        on:click={() => (selectedTab = "mortgage")}
        type="button"
    >
        💰
        <span><h1 class="text-xs font-bold leading-tight">Mortgage</h1></span>
    </button>
</div>

<style>
    .tabs {
        display: flex;
        justify-content: space-between;
        margin-bottom: 2rem;
        background: #121212;
        border-radius: 16px;
        padding: 0.5rem;
        gap: 0.5rem;
    }
    .tab {
        display: flex;
        flex-direction: column;
        align-items: center;
        cursor: pointer;
        padding: 0.2rem;
        flex: 1;
        border-radius: 12px;
        color: #ccc;
        transition:
            background 0.3s ease,
            color 0.3s ease;
    }
    .tab:hover {
        background: #2a2a2a;
    }
    .tab.active {
        background: rgba(186, 126, 254, 0.7);
        color: #fff;
    }
    .tab span {
        font-size: 0.9rem;
    }
</style>
