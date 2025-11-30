import grpc
from concurrent import futures
import time
import io
from PIL import Image
from transformers import pipeline # Hugging Face 的核心库

import moderation_pb2
import moderation_pb2_grpc

class ModerationService(moderation_pb2_grpc.ModerationServiceServicer):
    def __init__(self):
        print("正在加载 AI 模型 (Falconsai/nsfw_image_detection)... 第一次运行会自动下载，请耐心等待")
        # 初始化图像分类管道
        # 这个模型会返回 [{"label": "nsfw", "score": 0.99}, {"label": "normal", "score": 0.01}]
        self.classifier = pipeline("image-classification", model="Falconsai/nsfw_image_detection")
        print("AI 模型加载完毕，服务就绪！")

    def CheckImage(self, request, context):
        print(f"[Python] 收到图片审核请求: 文件名={request.file_name}, 大小={len(request.image_data)} bytes")
        
        try:
            # 1. 将 gRPC 传来的 bytes 转换为 PIL 图片对象
            image = Image.open(io.BytesIO(request.image_data))
            
            # 2. AI 推理
            # classifier 返回的是一个列表，包含所有类别的分数
            results = self.classifier(image)
            # 结果示例: [{'label': 'nsfw', 'score': 0.98}, {'label': 'normal', 'score': 0.02}]
            
            # 3. 解析结果
            is_nsfw = False
            score = 0.0
            
            # 找到 nsfw 标签的分数
            for res in results:
                if res['label'] == 'nsfw':
                    score = res['score']
                    break
            
            # 4. 判定逻辑：如果 NSFW 概率大于 80% (0.8)，则拦截
            threshold = 0.8
            if score > threshold:
                reason = f"AI 检测为敏感内容 (置信度: {score:.2%})"
                print(f"拦截违规图片: {request.file_name}, 原因: {reason}")
                return moderation_pb2.CheckResponse(is_safe=False, reason=reason)
            else:
                print(f"图片通过: {request.file_name} (NSFW 概率: {score:.2%})")
                return moderation_pb2.CheckResponse(is_safe=True, reason="Pass")

        except Exception as e:
            print(f"处理出错: {e}")
            # 如果图片损坏或无法识别，为了业务连续性，通常可以选择放行或报错
            # 这里演示报错
            return moderation_pb2.CheckResponse(is_safe=False, reason=f"图片处理失败: {str(e)}")

def serve():
    # 定义最大消息大小 (例如 50MB = 50 * 1024 * 1024)
    MAX_MESSAGE_LENGTH = 50 * 1024 * 1024

    # 建立 gRPC 服务器时，传入 options 参数
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10),
        options=[
            ('grpc.max_send_message_length', MAX_MESSAGE_LENGTH),
            ('grpc.max_receive_message_length', MAX_MESSAGE_LENGTH),
        ]
    )
    
    # 实例化服务类
    service_instance = ModerationService()
    
    moderation_pb2_grpc.add_ModerationServiceServicer_to_server(service_instance, server)
    
    server.add_insecure_port('[::]:50051')
    print(f"🚀 Python AI 内容审核微服务已启动 (Port: 50051) | 最大消息限制: {MAX_MESSAGE_LENGTH/1024/1024}MB...")
    try:
        server.start()
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()