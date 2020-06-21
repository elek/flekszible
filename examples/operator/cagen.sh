certstrap init --cn ca --passphrase ""
certstrap request-cert --common-name flekszible-operator --domain flekszible-operator,flekszible-operator.default,flekszible-operator.default.svc,flekszible-operator.default.svc.cluster.local
certstrap sign -CA ca flekszible-operator
cp out/flekszible-operator.crt configmaps/flekszible-operator-tls_server.crt
cp out/flekszible-operator.key configmaps/flekszible-operator-tls_server.key
#update resources/webhook-resigration with this:
cat out/ca.crt | base64 -w 0
