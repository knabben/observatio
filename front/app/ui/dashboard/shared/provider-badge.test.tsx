import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";

import {ProviderBadge} from "./provider-badge";
import {InfraCapabilityProvider} from "./infra-capability-context";
import {InfrastructureCapability} from "@/app/lib/data";

describe("ProviderBadge", () => {
  const capability: InfrastructureCapability = {
    docker: {installed: true, version: "v1.10.10"},
    vsphere: {installed: true, version: "v1.12.0"},
  };

  it("shows the Docker label with version", () => {
    render(
      <InfraCapabilityProvider value={capability}>
        <ProviderBadge provider="docker"/>
      </InfraCapabilityProvider>,
    );
    expect(screen.getByText("Docker v1.10.10")).toBeInTheDocument();
  });

  it("shows the vSphere label with version", () => {
    render(
      <InfraCapabilityProvider value={capability}>
        <ProviderBadge provider="vsphere"/>
      </InfraCapabilityProvider>,
    );
    expect(screen.getByText("vSphere v1.12.0")).toBeInTheDocument();
  });

  it("falls back to the bare label when no version is known", () => {
    render(
      <InfraCapabilityProvider value={{docker: {installed: true, version: ""}, vsphere: {installed: false, version: ""}}}>
        <ProviderBadge provider="docker"/>
      </InfraCapabilityProvider>,
    );
    expect(screen.getByText("Docker")).toBeInTheDocument();
  });

  it("shows Unknown for an unrecognized provider", () => {
    render(
      <InfraCapabilityProvider value={capability}>
        <ProviderBadge provider="aws"/>
      </InfraCapabilityProvider>,
    );
    expect(screen.getByText("Unknown")).toBeInTheDocument();
  });

  it("shows Unknown when no provider is given", () => {
    render(
      <InfraCapabilityProvider value={capability}>
        <ProviderBadge/>
      </InfraCapabilityProvider>,
    );
    expect(screen.getByText("Unknown")).toBeInTheDocument();
  });
});
