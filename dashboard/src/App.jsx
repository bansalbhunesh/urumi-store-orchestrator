import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Plus, LayoutGrid, Github, Sparkles } from 'lucide-react';
import StoreList from './components/StoreList';
import CreateStoreModal from './components/CreateStoreModal';

function App() {
    const [stores, setStores] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false);

    const fetchStores = async () => {
        try {
            const response = await axios.get('/api/stores');
            setStores(response.data);
        } catch (error) {
            console.error('Error fetching stores:', error);
            // Fallback for demo purposes if backend isn't reachable immediately
            if (stores.length === 0) setStores([]);
        }
    };

    useEffect(() => {
        fetchStores();
        const interval = setInterval(fetchStores, 5000); // Poll every 5s
        return () => clearInterval(interval);
    }, []);

    const handleCreateStore = async (storeData) => {
        try {
            await axios.post('/api/stores', storeData);
            fetchStores(); // Refresh list immediately
        } catch (error) {
            console.error('Error creating store:', error);
            alert('Failed to create store');
        }
    };

    const handleDeleteStore = async (id) => {
        if (!confirm('Are you sure you want to delete this store?')) return;
        try {
            await axios.delete(`/api/stores/${id}`);
            fetchStores();
        } catch (error) {
            console.error('Error deleting store:', error);
        }
    };

    return (
        <div className="min-h-screen selection:bg-violet-500/30">
            {/* Navbar */}
            <nav className="sticky top-0 z-40 w-full border-b border-white/5 bg-slate-950/60 backdrop-blur-xl supports-[backdrop-filter]:bg-slate-950/30">
                <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
                    <div className="flex items-center gap-3 group cursor-default">
                        <div className="relative">
                            <div className="absolute inset-0 bg-violet-500 blur-lg opacity-40 group-hover:opacity-60 transition-opacity"></div>
                            <div className="relative p-2.5 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-xl shadow-xl shadow-violet-500/20 group-hover:scale-105 transition-transform duration-300">
                                <LayoutGrid className="w-5 h-5 text-white" />
                            </div>
                        </div>
                        <div className="flex flex-col">
                            <span className="font-bold text-xl tracking-tight leading-none">Urumi</span>
                            <span className="text-xs font-medium text-violet-400 tracking-wider">ORCHESTRATOR</span>
                        </div>
                    </div>

                    <div className="flex items-center gap-4">
                        <a href="https://github.com/urumi-ai" target="_blank" className="p-2 text-slate-400 hover:text-white hover:bg-white/5 rounded-full transition-all">
                            <Github className="w-5 h-5" />
                        </a>
                    </div>
                </div>
            </nav>

            {/* Main Content */}
            <main className="max-w-7xl mx-auto px-6 py-12">
                <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6 mb-12 animate-fade-in-up">
                    <div>
                        <div className="admin-badge inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-violet-500/10 border border-violet-500/20 text-violet-300 text-xs font-semibold mb-3">
                            <Sparkles className="w-3 h-3" />
                            <span>V1.0.0 Stable</span>
                        </div>
                        <h1 className="text-4xl md:text-5xl font-bold mb-3 text-gradient">
                            Store Dashboard
                        </h1>
                        <p className="text-slate-400 text-lg max-w-xl leading-relaxed">
                            Monitor and manage your high-performance e-commerce deployments from a unified control plane.
                        </p>
                    </div>

                    <button
                        onClick={() => setIsModalOpen(true)}
                        className="group relative px-6 py-3.5 rounded-xl bg-white text-slate-950 font-bold overflow-hidden shadow-[0_0_20px_rgba(255,255,255,0.3)] hover:shadow-[0_0_30px_rgba(255,255,255,0.4)] transition-all active:scale-[0.98]"
                    >
                        <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/50 to-transparent translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-700"></div>
                        <span className="relative flex items-center gap-2">
                            <Plus className="w-5 h-5" />
                            Provision Store
                        </span>
                    </button>
                </div>

                <div className="animate-fade-in-up" style={{ animationDelay: '0.1s' }}>
                    <StoreList stores={stores} onDelete={handleDeleteStore} />
                </div>
            </main>

            <CreateStoreModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onCreate={handleCreateStore}
            />
        </div>
    );
}

export default App;
