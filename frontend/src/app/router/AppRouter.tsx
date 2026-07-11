import { lazy } from 'react';
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';

const LobbyPage = lazy(() => import('../../pages/lobby/LobbyPage'));
const RoomPage = lazy(() => import('../../pages/room/RoomPage'));
const MatchPage = lazy(() => import('../../pages/match/MatchPage'));

const router = createBrowserRouter([
  { path: '/', element: <LobbyPage /> },
  { path: '/rooms/:roomId', element: <RoomPage /> },
  { path: '/matches/:matchId', element: <MatchPage /> },
  { path: '*', element: <Navigate to="/" replace /> },
]);

export function AppRouter() {
  return <RouterProvider router={router} />;
}
