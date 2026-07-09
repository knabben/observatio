import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";
import KubeadmControlPlaneTable from "./table";
import {KubeadmControlPlaneType} from "./types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

describe("KubeadmControlPlaneTable", () => {
  it("renders a partial KubeadmControlPlane (missing status) without throwing", () => {
    const partial = [{metadata: {name: "kcp1"}}] as KubeadmControlPlaneType[];
    wrap(<KubeadmControlPlaneTable kcps={partial} select={() => {}}/>);
    expect(screen.getByRole("button", {name: /select kcp1/i})).toBeInTheDocument();
  });

  it("shows unknown readiness as 'Unknown', not failed", () => {
    const partial = [{metadata: {name: "kcp1"}}] as KubeadmControlPlaneType[];
    wrap(<KubeadmControlPlaneTable kcps={partial} select={() => {}}/>);
    expect(screen.getByRole("img", {name: "Unknown"})).toBeInTheDocument();
    expect(screen.queryByRole("img", {name: "Not ready"})).not.toBeInTheDocument();
  });

  it("renders a labeled empty state for an empty collection without crashing", () => {
    wrap(<KubeadmControlPlaneTable kcps={[]} select={() => {}}/>);
    expect(screen.getByText(/no kubeadm control planes found/i)).toBeInTheDocument();
  });
});
