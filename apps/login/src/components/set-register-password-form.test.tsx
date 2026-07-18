import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { create } from "@zitadel/client";
import { PasswordComplexitySettingsSchema } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { afterEach, describe, expect, test, vi } from "vitest";
import { registerUser } from "@/lib/server/register";
import { SetRegisterPasswordForm } from "./set-register-password-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/register", () => ({
  registerUser: vi.fn(),
}));

vi.mock("@/lib/client", () => ({
  handleServerActionResponse: vi.fn(),
}));

const defaultComplexitySettings = create(PasswordComplexitySettingsSchema, {
  minLength: 8n,
  requiresUppercase: false,
  requiresLowercase: false,
  requiresNumber: false,
  requiresSymbol: false,
});

describe("SetRegisterPasswordForm", () => {
  afterEach(cleanup);

  test("should autofocus the password input on mount", () => {
    const { getByTestId } = render(
      <SetRegisterPasswordForm
        passwordComplexitySettings={defaultComplexitySettings}
        email="test@example.com"
        firstname="Test"
        lastname="User"
        organization="org-1"
      />,
    );
    expect(getByTestId("password-text-input")).toHaveFocus();
  });

  test("should carry the username and display name chosen on the previous step through to registerUser", async () => {
    const { getByTestId } = render(
      <SetRegisterPasswordForm
        passwordComplexitySettings={defaultComplexitySettings}
        email="test@example.com"
        firstname="Test"
        lastname="User"
        username="chosen-username"
        displayname="Chosen Display Name"
        organization="org-1"
      />,
    );

    fireEvent.change(getByTestId("password-text-input"), { target: { value: "Password1!" } });
    fireEvent.change(getByTestId("password-confirm-text-input"), { target: { value: "Password1!" } });
    fireEvent.click(getByTestId("submit-button"));

    await waitFor(() =>
      expect(registerUser).toHaveBeenCalledWith(
        expect.objectContaining({
          username: "chosen-username",
          displayName: "Chosen Display Name",
        }),
      ),
    );
  });
});
