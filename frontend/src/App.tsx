import { BrowserRouter, Routes, Route } from "react-router";
import { ProtectedRoute } from "@/components/ui/ProtectedRoute";
import { PublicRoute } from "@/components/ui/PublicRoute";
import { AppLayout } from "@/components/ui/AppLayout";
import { DashboardPage } from "@/pages/DashboardPage";
import { LoginPage } from "@/pages/LoginPage";
import { RegisterPage } from "@/pages/RegisterPage";
import { NotFoundPage } from "@/pages/NotFoundPage";
import { NewProjectPage } from "@/pages/NewProjectPage";
import { ProjectDetailPage } from "@/pages/ProjectDetailPage";
import { DeploymentDetailPage } from "@/pages/DeploymentDetailPage";
import { LandingPage } from "@/pages/LandingPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Marketing landing — public */}
        <Route path="/" element={<LandingPage />} />

        {/* Auth pages — redirect to dashboard if logged in */}
        <Route element={<PublicRoute />}>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
        </Route>

        {/* Protected dashboard pages */}
        <Route element={<ProtectedRoute />}>
          <Route element={<AppLayout />}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/projects/new" element={<NewProjectPage />} />
            <Route path="/projects/:id" element={<ProjectDetailPage />} />
            <Route
              path="/projects/:projectId/deployments/:deploymentId"
              element={<DeploymentDetailPage />}
            />
          </Route>
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
