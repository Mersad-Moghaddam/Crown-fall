import { Component, type ErrorInfo, type ReactNode, Suspense } from 'react';
import { AppProviders } from './providers/AppProviders';
import { AppRouter } from './router/AppRouter';
import { LoadingScreen } from '../shared/ui/LoadingScreen';

type BoundaryState = { failed: boolean };

class ErrorBoundary extends Component<{ children: ReactNode }, BoundaryState> {
  state: BoundaryState = { failed: false };

  static getDerivedStateFromError(): BoundaryState {
    return { failed: true };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('Application boundary caught an error', error.name, info.componentStack);
  }

  render() {
    if (this.state.failed) {
      return (
        <main className="centered">
          <h1>The Chronicle was interrupted</h1>
          <p>Reload to reconnect safely.</p>
        </main>
      );
    }
    return this.props.children;
  }
}

export function App() {
  return (
    <ErrorBoundary>
      <AppProviders>
        <Suspense fallback={<LoadingScreen />}>
          <AppRouter />
        </Suspense>
      </AppProviders>
    </ErrorBoundary>
  );
}
