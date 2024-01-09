exports.handler = async (event, context) => {
  for (let i = 0; i < 5000; i++) {
    console.log("@".repeat(100));
  }
  const response = {
    statusCode: 200,
    body: "hello, world",
  };
  return response;
};
