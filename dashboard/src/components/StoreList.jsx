import React from 'react';
import { ExternalLink, Trash2, Box, RefreshCw, ShoppingCart } from 'lucide-react';

const StoreStatus = ({ status }) => {
    const styles = {
        Provisioning: "bg-amber-500/10 text-amber-500 border-amber-500/20",
        Ready: "bg-emerald-500/10 text-emerald-500 border-emerald-500/20",
        Failed: "bg-red-500/10 text-red-500 border-red-500/20",
    };

    const statusStyle = styles[status] || "bg-slate-500/10 text-slate-500";

    return (
        <div className={`px-2 py-1 rounded-full text-xs font-medium border flex items-center gap-1.5 ${statusStyle}`}>
            {status === 'Provisioning' && <RefreshCw className="w-3 h-3 animate-spin" />}
            {status}
        </div>
    );
};

const StoreCard = ({ store, onDelete }) => {
    return (
        <div className="glass-card rounded-xl p-5 group flex flex-col h-full border hover:border-violet-500/30">
            <div className="flex justify-between items-start mb-4">
                <div className="p-2.5 rounded-lg bg-gradient-to-br from-violet-500/20 to-indigo-500/10 border border-violet-500/10">
                    <ShoppingCart className="w-6 h-6 text-violet-400" />
                </div>
                <StoreStatus status={store.status} />
            </div>

            <h3 className="text-lg font-semibold text-white mb-1 group-hover:text-violet-300 transition-colors">
                {store.name}
            </h3>
            <p className="text-sm text-slate-400 mb-6 font-mono opacity-80">
                {store.type}
            </p>

            <div className="mt-auto flex flex-col gap-3">
                {store.status === 'Ready' && (
                    <a
                        href={store.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center justify-center gap-2 w-full py-2 rounded-lg bg-white/5 hover:bg-white/10 text-sm font-medium transition-colors border border-white/5"
                    >
                        Visit Store <ExternalLink className="w-4 h-4" />
                    </a>
                )}

                <button
                    onClick={() => onDelete(store.id)}
                    className="flex items-center justify-center gap-2 w-full py-2 rounded-lg text-red-400/80 hover:bg-red-500/10 hover:text-red-400 text-sm font-medium transition-all"
                >
                    Delete <Trash2 className="w-4 h-4" />
                </button>
            </div>
        </div>
    );
};

export default function StoreList({ stores, onDelete }) {
    if (stores.length === 0) {
        return (
            <div className="text-center py-20 opacity-50 border-2 border-dashed border-slate-800 rounded-2xl">
                <Box className="w-12 h-12 mx-auto mb-4 text-slate-600" />
                <p className="text-lg text-slate-400">No stores found. Create one to get started.</p>
            </div>
        );
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {stores.map(store => (
                <StoreCard key={store.id} store={store} onDelete={onDelete} />
            ))}
        </div>
    );
}
