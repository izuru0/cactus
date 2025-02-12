import { type } from "os";

export type SupportedFunctions =
  | "registerHistoryData"
  | "getLineage"
  | "searchByHeader"
  | "searchByGlobalData"
  | "status";

export type AuthInfoAccessTokenArgsType = {
  accessToken: string;
  trustAgentId: string;
};

export type AuthInfoSubscriptionKeyArgsType = {
  subscriptionKey: string;
  trustAgentId: string;
  trustAgentRole: string;
  trustUserId: string;
  trustUserRole: string;
};

export type AuthInfoArgsType =
  | AuthInfoAccessTokenArgsType
  | AuthInfoSubscriptionKeyArgsType;

export function isAuthInfoAccessTokenArgsType(
  authInfo: AuthInfoArgsType,
): authInfo is AuthInfoAccessTokenArgsType {
  const typedAuthInfo = authInfo as AuthInfoAccessTokenArgsType;
  return (
    typedAuthInfo &&
    typeof typedAuthInfo.accessToken !== "undefined" &&
    typeof typedAuthInfo.trustAgentId !== "undefined"
  );
}

export function isAuthInfoSubscriptionKeyArgsType(
  authInfo: AuthInfoArgsType,
): authInfo is AuthInfoSubscriptionKeyArgsType {
  const typedAuthInfo = authInfo as AuthInfoSubscriptionKeyArgsType;
  return (
    typedAuthInfo &&
    typeof typedAuthInfo.subscriptionKey !== "undefined" &&
    typeof typedAuthInfo.trustAgentId !== "undefined" &&
    typeof typedAuthInfo.trustAgentRole !== "undefined" &&
    typeof typedAuthInfo.trustUserId !== "undefined" &&
    typeof typedAuthInfo.trustUserRole !== "undefined"
  );
}

export type FunctionArgsType = {
  method: {
    type: SupportedFunctions;
    authInfo: AuthInfoArgsType;
  };
  args: any;
  reqID?: string;
};
