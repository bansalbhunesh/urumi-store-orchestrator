import React from 'react';
import { ExternalLink, Trash2, Box, RefreshCw, ShoppingCart, Globe, Clock } from 'lucide-react';

const StoreStatus = ({ status }) => {
    const styles = {
        Provisioning: "bg-amber-500/10 text-amber-500 border-amber-500/20 shadow-[0_0_10px_rgba(245,158,11,0.1)]",
        Ready: "bg-emerald-500/10 text-emerald-500 border-emerald-500/20 shadow-[0_0_10px_rgba(16,185,129,0.1)]",
        Failed: "bg-red-500/10 text-red-500 border-red-500/20 shadow-[0_0_10px_rgba(239,68,68,0.1)]",
    };

    const statusStyle = styles[status] || "bg-slate-500/10 text-slate-500";

    return (
        <div className={`px-2.5 py-1 rounded-full text-xs font-semibold border flex items-center gap-1.5 ${statusStyle}`}>
            {status === 'Provisioning' && <RefreshCw className="w-3 h-3 animate-spin" />}
            {status === 'Ready' && <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />}
            {status}
        </div>
    );
};

const StoreCard = ({ store, onDelete }) => {
    return (
        <div className="glass-card rounded-2xl p-6 group flex flex-col h-full relative overflow-hidden">
            {/* Gradient Orb */}
            <div className="absolute top-0 right-0 w-32 h-32 bg-violet-500/10 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2 group-hover:bg-violet-500/20 transition-all duration-500"></div>

            <div className="flex justify-between items-start mb-6 relative z-10">
                <div className="p-3 rounded-xl bg-gradient-to-br from-violet-500/20 to-indigo-500/10 border border-violet-500/10 shadow-lg shadow-violet-500/5 group-hover:scale-110 transition-transform duration-300">
                    <ShoppingCart className="w-6 h-6 text-violet-300" />
                </div>
                <StoreStatus status={store.status} />
            </div>

            <div className="relative z-10 mb-6">
                <h3 className="text-xl font-bold text-white mb-1 group-hover:text-violet-300 transition-colors tracking-tight">
                    {store.name}
                </h3>
                <div className="flex items-center gap-4 text-xs font-medium text-slate-500 mt-3">
                    <span className="flex items-center gap-1 bg-white/5 px-2 py-1 rounded-md border border-white/5">
                        <Globe className="w-3 h-3" /> {store.type}
                    </span>
                    <span className="flex items-center gap-1" title={new Date(store.created_at).toLocaleString()}>
                        <Clock className="w-3 h-3" />
                        {(() => {
                            const date = new Date(store.created_at);
                            const now = new Date();
                            const diffInSeconds = Math.floor((now - date) / 1000);

                            if (diffInSeconds < 60) return 'Just now';
                            if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
                            if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`;
                            return `${Math.floor(diffInSeconds / 86400)}d ago`;
                        })()}
                    </span>
                </div>
            </div>

            <div className="mt-auto flex flex-col gap-3 relative z-10">
                {store.status === 'Ready' && (
                    <a
                        href={store.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center justify-center gap-2 w-full py-2.5 rounded-xl bg-white text-slate-950 font-bold hover:bg-slate-200 transition-all active:scale-[0.98] shadow-lg shadow-white/5"
                    >
                        Visit Store <ExternalLink className="w-4 h-4" />
                    </a>
                )}

                <button
                    onClick={() => onDelete(store.id)}
                    className="flex items-center justify-center gap-2 w-full py-2.5 rounded-xl text-slate-400 hover:text-red-400 hover:bg-red-500/10 transition-all text-sm font-medium border border-transparent hover:border-red-500/10"
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
            <div className="text-center py-24 glass-panel rounded-3xl border-dashed border-slate-700/50">
                <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-slate-800/50 flex items-center justify-center animate-pulse">
                    <Box className="w-10 h-10 text-slate-600" />
                </div>
                <h3 className="text-xl font-bold text-white mb-2">No active stores</h3>
                <p className="text-slate-400 max-w-sm mx-auto">
                    You haven't provisioned any stores yet. Click the "Provision Store" button to get started.
                </p>
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
