# Grpc-Demo
Just a simple example to demonstrate the four ways of grpc.


以导航的形式，针对四种grpc的传输形式，客户端与服务端的样例 定义四种不同的信息类型：Point,Rectangle,Feature,RouteSummary以及chat 定义四个方法

1. GetFeatures (输入为一个Point, 返回这个点的Feature)
2. ListFeatures (输入一个为Rectangle, 输出流这个区域所有的Feature)
3. RecordRoute (输入流为每个时间点的位置Point, 返回一个RouteSummary)
4. Recommend (输入流RecommendationRequest, 输出流Feature)

