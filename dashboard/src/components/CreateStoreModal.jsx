import React, { useState } from 'react';
import { X, Sparkles, Server } from 'lucide-react';

export default function CreateStoreModal({ isOpen, onClose, onCreate, isLoading }) {
    const [name, setName] = useState('');
    const [type, setType] = useState('woocommerce');
    const [isSubmitting, setIsSubmitting] = useState(false);

    if (!isOpen) return null;

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (isSubmitting) return;
        
        setIsSubmitting(true);
        try {
            await onCreate({ name, type });
            onClose();
            setName('');
            setType('woocommerce');
        } catch (error) {
            // Error is handled by parent component
            console.error('Create store error:', error);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-slate-950/60 backdrop-blur-md transition-all duration-300"
                onClick={onClose}
            ></div>

            {/* Modal */}
            <div className="relative w-full max-w-md glass-panel rounded-2xl overflow-hidden animate-fade-in-up">
                <div className="p-6 border-b border-white/5 flex justify-between items-center bg-white/5">
                    <h2 className="text-xl font-bold flex items-center gap-2 text-white">
                        <Sparkles className="w-5 h-5 text-violet-400" />
                        Provision New Store
                    </h2>
                    <button onClick={onClose} className="p-2 rounded-full hover:bg-white/10 text-slate-400 hover:text-white transition-colors">
                        <X className="w-5 h-5" />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="p-6 space-y-6">
                    <div className="space-y-2">
                        <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Store Name</label>
                        <input
                            type="text"
                            required
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            disabled={isSubmitting}
                            className="w-full bg-slate-950/50 border border-slate-700/50 rounded-xl px-4 py-3.5 text-white focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-violet-500 transition-all placeholder-slate-600 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                            placeholder="e.g. My Awesome Shop"
                            minLength={2}
                            maxLength={50}
                        />
                    </div>

                    <div className="space-y-4">
                        <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Store Engine</label>
                        <div className="grid grid-cols-2 gap-4">
                            <button
                                type="button"
                                disabled={isSubmitting}
                                className={`p-4 rounded-xl border flex flex-col items-center gap-3 transition-all duration-200 ${type === 'woocommerce' ? 'bg-violet-500/20 border-violet-500/50 text-white shadow-[0_0_15px_rgba(139,92,246,0.15)] scale-[1.02]' : 'bg-slate-900/40 border-slate-700/30 text-slate-400 hover:bg-slate-800/60 hover:border-slate-600'} disabled:opacity-50 disabled:cursor-not-allowed`}
                                onClick={() => setType('woocommerce')}
                            >
                                <div className={`w-10 h-10 rounded-full flex items-center justify-center ${type === 'woocommerce' ? 'bg-violet-500' : 'bg-slate-800'}`}>
                                    <span className="text-lg font-bold">W</span>
                                </div>
                                <span className="font-semibold text-sm">WooCommerce</span>
                            </button>

                            <button
                                type="button"
                                disabled={isSubmitting}
                                className={`p-4 rounded-xl border flex flex-col items-center gap-3 transition-all duration-200 ${type === 'medusa' ? 'bg-violet-500/20 border-violet-500/50 text-white shadow-[0_0_15px_rgba(139,92,246,0.15)] scale-[1.02]' : 'bg-slate-900/40 border-slate-700/30 text-slate-400 hover:bg-slate-800/60 hover:border-slate-600'} disabled:opacity-50 disabled:cursor-not-allowed`}
                                onClick={() => setType('medusa')}
                            >
                                <div className={`w-10 h-10 rounded-full flex items-center justify-center ${type === 'medusa' ? 'bg-violet-500' : 'bg-slate-800'}`}>
                                    <span className="text-lg font-bold">M</span>
                                </div>
                                <span className="font-semibold text-sm">MedusaJS</span>
                            </button>
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={isSubmitting || isLoading || !name.trim()}
                        className="w-full py-4 rounded-xl bg-white hover:bg-slate-200 text-slate-950 font-bold shadow-lg shadow-white/5 transition-all active:scale-[0.98] disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2 mt-2"
                    >
                        {(isSubmitting || isLoading) ? <Server className="w-4 h-4 animate-spin" /> : 'Launch Store'}
                    </button>
                </form>
            </div>
        </div>
    );
}
