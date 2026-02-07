import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Plus, LayoutGrid, Github } from 'lucide-react';
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
        <div className="min-h-screen bg-slate-950 text-white selection:bg-violet-500/30">
            {/* Navbar */}
            <nav className="sticky top-0 z-40 w-full border-b border-slate-800 bg-slate-950/80 backdrop-blur-xl">
                <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <div className="p-2 bg-violet-600 rounded-lg shadow-lg shadow-violet-500/20">
                            <LayoutGrid className="w-5 h-5 text-white" />
                        </div>
                        <span className="font-bold text-xl tracking-tight">Urumi<span className="text-violet-400">Stores</span></span>
                    </div>

                    <div className="flex items-center gap-4">
                        <a href="https://github.com/urumi-ai" target="_blank" className="text-slate-400 hover:text-white transition-colors">
                            <Github className="w-5 h-5" />
                        </a>
                    </div>
                </div>
            </nav>

            {/* Main Content */}
            <main className="max-w-7xl mx-auto px-6 py-12">
                <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-12">
                    <div>
                        <h1 className="text-4xl font-bold mb-2 bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
                            Your Stores
                        </h1>
                        <p className="text-slate-400 text-lg">
                            Manage your ecommerce deployments
                        </p>
                    </div>

                    <button
                        onClick={() => setIsModalOpen(true)}
                        className="px-6 py-3 rounded-xl bg-white text-slate-950 font-bold hover:bg-slate-200 transition-colors shadow-xl shadow-white/5 flex items-center gap-2 active:scale-95"
                    >
                        <Plus className="w-5 h-5" />
                        New Store
                    </button>
                </div>

                <StoreList stores={stores} onDelete={handleDeleteStore} />
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
