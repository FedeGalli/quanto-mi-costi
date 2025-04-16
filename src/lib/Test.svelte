<script>
    export let sidebarItems = [
        { name: "Menu Item 1", href: "#" },
        { name: "Menu Item 2", href: "#" },
        { name: "Menu Item 3", href: "#" },
    ];
</script>

<div class="h-screen bg-gray-100 overflow-hidden">
    <div class="flex h-full">
        <!-- Left sidebar -->
        <div class="w-64 bg-gray-800 text-white p-4 flex flex-col">
            <h1 class="text-xl font-bold mb-6">Sidebar</h1>
            <nav class="flex-1">
                <ul class="space-y-2">
                    {#each sidebarItems as item}
                        <li>
                            <a
                                href={item.href}
                                class="block px-4 py-2 rounded hover:bg-gray-700"
                            >
                                {item.name}
                            </a>
                        </li>
                    {/each}
                </ul>
            </nav>
            <div class="mt-auto pt-4 border-t border-gray-700">
                <p class="text-sm text-gray-400">Footer content</p>
            </div>
        </div>

        <!-- Main content -->
        <div class="flex-1 flex flex-col overflow-hidden">
            <header class="bg-white shadow-sm p-4">
                <h2 class="text-xl font-semibold">Main Content</h2>
            </header>

            <h1>CIAONEEE</h1>

            <main class="flex-1 overflow-y-auto p-4 bg-white">
                <div class="max-w-4xl mx-auto">
                    <slot></slot>
                </div>
            </main>
        </div>
    </div>
</div>
