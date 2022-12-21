// https://blakewilliams.me/posts/handling-macos-url-schemes-with-go

#import <Cocoa/Cocoa.h>

extern void HandleURL(char*);

@interface CallbackHandlerAppDelegate: NSObject<NSApplicationDelegate>
    - (void)handleGetURLEvent:(NSAppleEventDescriptor *) event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void RunApp();