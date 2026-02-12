import React from 'react';
import { AlertTriangle, RefreshCw } from 'lucide-react';

class ErrorBoundary extends React.Component {
    constructor(props) {
        super(props);
        this.state = { hasError: false, error: null, errorInfo: null };
    }

    static getDerivedStateFromError(error) {
        return { hasError: true };
    }

    componentDidCatch(error, errorInfo) {
        this.setState({
            error: error,
            errorInfo: errorInfo
        });
        console.error('ErrorBoundary caught an error:', error, errorInfo);
    }

    render() {
        if (this.state.hasError) {
            return (
                <div className="min-h-screen selection:bg-violet-500/30 flex items-center justify-center p-6">
                    <div className="glass-panel rounded-2xl p-8 max-w-md w-full text-center">
                        <div className="w-16 h-16 mx-auto mb-6 rounded-full bg-red-500/10 flex items-center justify-center">
                            <AlertTriangle className="w-8 h-8 text-red-500" />
                        </div>
                        <h1 className="text-2xl font-bold text-white mb-4">Something went wrong</h1>
                        <p className="text-slate-400 mb-6">
                            The dashboard encountered an unexpected error. Please try refreshing the page.
                        </p>
                        <button
                            onClick={() => window.location.reload()}
                            className="w-full py-3 rounded-xl bg-white hover:bg-slate-200 text-slate-950 font-bold shadow-lg shadow-white/5 transition-all active:scale-[0.98] flex items-center justify-center gap-2"
                        >
                            <RefreshCw className="w-4 h-4" />
                            Refresh Page
                        </button>
                        {process.env.NODE_ENV === 'development' && this.state.error && (
                            <details className="mt-6 text-left">
                                <summary className="cursor-pointer text-sm text-slate-500 hover:text-slate-400">
                                    Error Details
                                </summary>
                                <pre className="mt-2 p-4 bg-slate-900/50 rounded-lg text-xs text-red-400 overflow-auto">
                                    {this.state.error && this.state.error.toString()}
                                    <br />
                                    {this.state.errorInfo.componentStack}
                                </pre>
                            </details>
                        )}
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}

export default ErrorBoundary;
