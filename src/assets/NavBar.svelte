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
        <span>Summary</span>
    </button>
    <button
        class="tab {selectedTab === 'base' ? 'active' : ''}"
        on:click={() => (selectedTab = "base")}
        type="button"
    >
        📦
        <span>Base Costs</span>
    </button>

    <button
        class="tab {selectedTab === 'mortgage' ? 'active' : ''}"
        on:click={() => (selectedTab = "mortgage")}
        type="button"
    >
        💰
        <span>Mortgage</span>
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
        padding: 0.5rem;
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
        background: rgba(54, 65, 83, 1);
        color: #fff;
    }
    .tab span {
        margin-top: 0.2rem;
        font-size: 0.9rem;
    }
</style>
