#include <gphoto2/gphoto2.h>

extern void wrapperInfoCallback(char* p0);
extern void wrapperErrorCallback(char* p0);
extern void wrapperLoggingCallback(int  logLevel, char* domain, char* data);


void
ctx_error_func (GPContext *context, const char *str, void *data)
{
	wrapperErrorCallback((char*)str);
}

void
ctx_status_func (GPContext *context, const char *str, void *data)
{
	wrapperInfoCallback((char*)str);
}


void 
loger_func(GPLogLevel level, const char *domain, const char *str, void *data) 
{
	wrapperLoggingCallback((int)level, (char*) domain, (char*) str);
}