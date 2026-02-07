import React, { useState } from 'react';
import { X, Sparkles, Server } from 'lucide-react';

export default function CreateStoreModal({ isOpen, onClose, onCreate }) {
    const [name, setName] = useState('');
    const [type, setType] = useState('woocommerce');
    const [loading, setLoading] = useState(false);

    if (!isOpen) return null;

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        await onCreate({ name, type });
        setLoading(false);
        onClose();
        setName('');
        setType('woocommerce');
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-slate-950/80 backdrop-blur-sm"
                onClick={onClose}
            ></div>

            {/* Modal */}
            <div className="relative w-full max-w-md bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-200">
                <div className="p-6 border-b border-slate-800 flex justify-between items-center bg-slate-900/50">
                    <h2 className="text-xl font-semibold flex items-center gap-2">
                        <Sparkles className="w-5 h-5 text-violet-400" />
                        Provision New Store
                    </h2>
                    <button onClick={onClose} className="text-slate-500 hover:text-white transition-colors">
                        <X className="w-5 h-5" />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="p-6 space-y-6">
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-slate-400">Store Name</label>
                        <input
                            type="text"
                            required
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500 transition-all placeholder-slate-600"
                            placeholder="e.g. My Awesome Shop"
                        />
                    </div>

                    <div className="space-y-4">
                        <label className="text-sm font-medium text-slate-400">Store Engine</label>
                        <div className="grid grid-cols-2 gap-4">
                            <button
                                type="button"
                                className={`p-4 rounded-xl border flex flex-col items-center gap-2 transition-all ${type === 'woocommerce' ? 'bg-violet-500/10 border-violet-500 text-violet-300' : 'bg-slate-950 border-slate-800 text-slate-500 hover:border-slate-700'}`}
                                onClick={() => setType('woocommerce')}
                            >
                                <div className="w-8 h-8 rounded-full bg-current opacity-20"></div>
                                <span className="font-medium">WooCommerce</span>
                            </button>

                            <button
                                type="button"
                                className={`p-4 rounded-xl border flex flex-col items-center gap-2 transition-all ${type === 'medusa' ? 'bg-violet-500/10 border-violet-500 text-violet-300' : 'bg-slate-950 border-slate-800 text-slate-500 hover:border-slate-700'}`}
                                onClick={() => setType('medusa')}
                            >
                                <div className="w-8 h-8 rounded-full bg-current opacity-20"></div>
                                <span className="font-medium">MedusaJS</span>
                            </button>
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full py-3.5 rounded-xl bg-violet-600 hover:bg-violet-500 text-white font-semibold shadow-lg shadow-violet-500/20 transition-all active:scale-[0.98] disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                        {loading ? <Server className="w-4 h-4 animate-spin" /> : 'Launch Store'}
                    </button>
                </form>
            </div>
        </div>
    );
}
