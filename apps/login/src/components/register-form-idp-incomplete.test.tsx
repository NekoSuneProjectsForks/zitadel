import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { RegisterFormIDPIncomplete } from "./register-form-idp-incomplete";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/register", () => ({
  registerUserAndLinkToIDP: vi.fn(),
}));

vi.mock("@/lib/client", () => ({
  handleServerActionResponse: vi.fn(),
}));

const defaultProps = {
  organization: "org-1",
  idpIntent: { idpIntentId: "intent-1", idpIntentToken: "token-1" },
  idpUserId: "user-1",
  idpId: "idp-1",
};

describe("RegisterFormIDPIncomplete", () => {
  afterEach(cleanup);

  test("should always show and autofocus an editable username input, even when the IDP suggested one", () => {
    const { getByTestId } = render(<RegisterFormIDPIncomplete {...defaultProps} />);
    expect(getByTestId("username-text-input")).toHaveFocus();
  });

  test("should still show and autofocus the username input when idpUserName is provided, pre-filled with the suggestion", () => {
    const { getByTestId } = render(<RegisterFormIDPIncomplete {...defaultProps} idpUserName="existing-user" />);
    const usernameInput = getByTestId("username-text-input") as HTMLInputElement;
    expect(usernameInput).toHaveFocus();
    expect(usernameInput.value).toBe("existing-user");
  });

  test("should always show an editable display name input", () => {
    const { getByTestId } = render(<RegisterFormIDPIncomplete {...defaultProps} idpUserName="existing-user" />);
    expect(getByTestId("displayname-text-input")).toBeInTheDocument();
  });
});
