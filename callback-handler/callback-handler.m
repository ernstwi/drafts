// https://blakewilliams.me/posts/handling-macos-url-schemes-with-go

#include "callback-handler.h"

@implementation CallbackHandlerAppDelegate
- (void)applicationWillFinishLaunching:(NSNotification *)aNotification
{
    NSAppleEventManager *appleEventManager = [NSAppleEventManager sharedAppleEventManager];
    [appleEventManager setEventHandler:self
                       andSelector:@selector(handleGetURLEvent:withReplyEvent:)
                       forEventClass:kInternetEventClass andEventID:kAEGetURL];
}

- (void)handleGetURLEvent:(NSAppleEventDescriptor *)event
           withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
    HandleURL((char*)[[[event paramDescriptorForKeyword:keyDirectObject] stringValue] UTF8String]);
}
@end

void RunApp() {
    [NSAutoreleasePool new];
    [NSApplication sharedApplication];
    CallbackHandlerAppDelegate *app = [CallbackHandlerAppDelegate alloc];
    [NSApp setDelegate:app];
    [NSApp run];
  }